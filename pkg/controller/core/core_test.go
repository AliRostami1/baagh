package core_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	// "github.com/AliRostami1/baagh/pkg/mockup"
	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"github.com/warthog618/gpiod/mockup"
	"github.com/warthog618/gpiod/uapi"
	"golang.org/x/sys/unix"
)

var platform Platform

func TestMain(m *testing.M) {
	var pname string
	flag.StringVar(&pname, "platform", "mockup", "test platform")
	flag.Parse()
	p, err := newPlatform(pname)
	if err != nil {
		fmt.Println("Platform not supported -", err)
		os.Exit(-1)
	}
	platform = p
	rc := m.Run()
	platform.Close()
	os.Exit(rc)
}

var (
	biasKernel               = mockup.Semver{5, 5}  // bias flags added
	setConfigKernel          = mockup.Semver{5, 5}  // setLineConfig ioctl added
	infoWatchKernel          = mockup.Semver{5, 7}  // watchLineInfo ioctl added
	uapiV2Kernel             = mockup.Semver{5, 10} // uapi v2 added
	eventClockRealtimeKernel = mockup.Semver{5, 11} // realtime event clock option added
)

func TestRequestItem(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// bad Chip name
	i, err := core.RequestItem(platform.Devpath()+"not", 1)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// negative Item offset
	i, err = core.RequestItem(platform.Devpath(), -1)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// out-of-range Item offset
	i, err = core.RequestItem(platform.Devpath(), platform.Lines())
	assert.NotNil(t, err)
	require.Nil(t, i)

	// success
	i, err = core.RequestItem(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, i)
	i.Close()
}

type gpiochip struct {
	name    string
	label   string
	devpath string
	lines   int
	// line triggered by TriggerIntr.
	intro     int
	introName string
	outo      int
	// floating lines - can be harmlessly set to outputs.
	ff []int
}

func (c *gpiochip) Name() string {
	return c.name
}

func (c *gpiochip) Label() string {
	return c.label
}
func (c *gpiochip) Devpath() string {
	return c.devpath
}

func (c *gpiochip) Lines() int {
	return c.lines
}

func (c *gpiochip) IntrLine() int {
	return c.intro
}

func (c *gpiochip) IntrName() string {
	return c.introName
}

func (c *gpiochip) OutLine() int {
	return c.outo
}

func (c *gpiochip) FloatingLines() []int {
	return c.ff
}

// two flavours of chip, raspberry and mockup.
type Platform interface {
	Name() string
	Label() string
	Devpath() string
	Lines() int
	IntrLine() int
	IntrName() string
	OutLine() int
	FloatingLines() []int
	TriggerIntr(int)
	ReadOut() int
	SupportsAsIs() bool
	Close()
}

type RaspberryPi struct {
	gpiochip
	chip  *gpiod.Chip
	wline *gpiod.Line
}

func isPi(path string) error {
	if err := gpiod.IsChip(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, unix.O_CLOEXEC, unix.O_RDONLY)
	if err != nil {
		return err
	}
	defer f.Close()
	ci, err := uapi.GetChipInfo(f.Fd())
	if err != nil {
		return err
	}
	label := uapi.BytesToString(ci.Label[:])
	if label != "pinctrl-bcm2835" && label != "pinctrl-bcm2711" {
		return fmt.Errorf("unsupported gpiochip - %s", label)
	}
	return nil
}

func newPi(path string) (*RaspberryPi, error) {
	if err := isPi(path); err != nil {
		return nil, err
	}
	ch, err := gpiod.NewChip(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			ch.Close()
		}
	}()
	pi := RaspberryPi{
		gpiochip: gpiochip{
			name:      "gpiochip0",
			label:     ch.Label,
			devpath:   path,
			lines:     int(ch.Lines()),
			intro:     rpi.J8p15,
			introName: "",
			outo:      rpi.J8p16,
			ff:        []int{rpi.J8p11, rpi.J8p12, rpi.J8p7, rpi.J8p13, rpi.J8p22},
		},
		chip: ch,
	}
	if ch.Label == "pinctrl-bcm2711" {
		pi.introName = "GPIO22"
	}
	// check J8p15 and J8p16 are tied
	w, err := ch.RequestLine(pi.outo, gpiod.AsOutput(1),
		gpiod.WithConsumer("gpiod-test-w"))
	if err != nil {
		return nil, err
	}
	defer w.Close()
	r, err := ch.RequestLine(pi.intro,
		gpiod.WithConsumer("gpiod-test-r"))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	v, _ := r.Value()
	if v != 1 {
		return nil, errors.New("J8p15 and J8p16 must be tied")
	}
	w.SetValue(0)
	v, _ = r.Value()
	if v != 0 {
		return nil, errors.New("J8p15 and J8p16 must be tied")
	}
	return &pi, nil
}

func (c *RaspberryPi) Close() {
	if c.wline != nil {
		c.wline.Close()
		c.wline = nil
	}
	// revert intr trigger line to input
	l, _ := c.chip.RequestLine(c.outo)
	l.Close()
	// revert floating lines to inputs
	ll, _ := c.chip.RequestLines(platform.FloatingLines())
	ll.Close()
	c.chip.Close()
}

func (c *RaspberryPi) OutLine() int {
	if c.wline != nil {
		c.wline.Close()
		c.wline = nil
	}
	return c.outo
}

func (c *RaspberryPi) ReadOut() int {
	r, err := c.chip.RequestLine(c.intro,
		gpiod.WithConsumer("gpiod-test-r"))
	if err != nil {
		return -1
	}
	defer r.Close()
	v, err := r.Value()
	if err != nil {
		return -1
	}
	return v
}

func (c *RaspberryPi) SupportsAsIs() bool {
	// RPi pinctrl-bcm2835 returns lines to input on release.
	return false
}

func (c *RaspberryPi) TriggerIntr(value int) {
	if c.wline != nil {
		c.wline.SetValue(value)
		return
	}
	w, _ := c.chip.RequestLine(c.outo, gpiod.AsOutput(value),
		gpiod.WithConsumer("gpiod-test-w"))
	c.wline = w
}

type Mockup struct {
	gpiochip
	m *mockup.Mockup
	c *mockup.Chip
}

func newMockup() (*Mockup, error) {
	m, err := mockup.New([]int{20}, true)
	if err != nil {
		return nil, err
	}
	c, err := m.Chip(0)
	if err != nil {
		return nil, err
	}
	return &Mockup{
		gpiochip{
			name:      c.Name,
			label:     c.Label,
			devpath:   c.DevPath,
			lines:     20,
			intro:     10,
			introName: "gpio-mockup-A-10",
			outo:      9,
			ff:        []int{11, 12, 15, 16, 9},
		}, m, c}, nil
}

func (c *Mockup) Close() {
	c.m.Close()
}

func (c *Mockup) ReadOut() int {
	v, err := c.c.Value(c.outo)
	if err != nil {
		return -1
	}
	return v
}

func (c *Mockup) SupportsAsIs() bool {
	return true
}

func (c *Mockup) TriggerIntr(value int) {
	c.c.SetValue(c.intro, value)
}

func newPlatform(pname string) (Platform, error) {
	switch pname {
	case "mockup":
		p, err := newMockup()
		if err != nil {
			return nil, fmt.Errorf("error loading gpio-mockup: %w", err)
		}
		return p, nil
	case "rpi":
		return newPi("/dev/gpiochip0")
	default:
		return nil, fmt.Errorf("unknown platform '%s'", pname)
	}
}

func requireKernel(t *testing.T, min mockup.Semver) {
	t.Helper()
	if err := mockup.CheckKernelVersion(min); err != nil {
		t.Skip(err)
	}
}

func requireABI(t *testing.T, chip *gpiod.Chip, abi int) {
	t.Helper()
	if chip.UapiAbiVersion() != abi {
		t.Skip(ErrorBadABIVersion{abi, chip.UapiAbiVersion()})
	}
}

// ErrorBadVersion indicates the kernel version is insufficient.
type ErrorBadABIVersion struct {
	Need int
	Have int
}

func (e ErrorBadABIVersion) Error() string {
	return fmt.Sprintf("require kernel ABI %d, but using %d", e.Need, e.Have)
}
