package core_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

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

// var (
// 	biasKernel               = mockup.Semver{5, 5}  // bias flags added
// 	setConfigKernel          = mockup.Semver{5, 5}  // setLineConfig ioctl added
// 	infoWatchKernel          = mockup.Semver{5, 7}  // watchLineInfo ioctl added
// 	uapiV2Kernel             = mockup.Semver{5, 10} // uapi v2 added
// 	eventClockRealtimeKernel = mockup.Semver{5, 11} // realtime event clock option added
// )

func TestRequestItem(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// fail: bad Chip name
	i, err := core.RequestItem(platform.Devpath()+"not", 1)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// fail: negative Item offset
	i, err = core.RequestItem(platform.Devpath(), -1)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// fail: out-of-range Item offset
	i, err = core.RequestItem(platform.Devpath(), platform.Lines())
	assert.NotNil(t, err)
	require.Nil(t, i)

	// success: without options
	i2, err := core.RequestItem(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, i2)
	err = i2.Close()
	assert.Nil(t, err)

	// success: with input(core.PullDisabled) option
	i3, err := core.RequestItem(platform.Devpath(), io, core.AsInput(core.PullDisabled))
	assert.Nil(t, err)
	require.NotNil(t, i3)
	err = i3.Close()
	assert.Nil(t, err)

	// success: with input(core.PullUp) option
	i4, err := core.RequestItem(platform.Devpath(), io, core.AsInput(core.PullUp))
	assert.Nil(t, err)
	require.NotNil(t, i4)
	err = i4.Close()
	assert.Nil(t, err)

	// success: with input(core.PullDown) option
	i5, err := core.RequestItem(platform.Devpath(), io, core.AsInput(core.PullDown))
	assert.Nil(t, err)
	require.NotNil(t, i5)
	err = i5.Close()
	assert.Nil(t, err)

	// success: with output(high)
	i6, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateActive))
	assert.Nil(t, err)
	require.NotNil(t, i6)
	err = i6.Close()
	assert.Nil(t, err)

	// success: with output(low)
	i7, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i7)
	err = i7.Close()
	assert.Nil(t, err)

	// success: multiple request with same config of the same line
	i8, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i8)
	i9, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i9)
	i10, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i10)

	assert.Equal(t, i8, i9)
	assert.Equal(t, i9, i10)

	err = i8.Close()
	assert.Nil(t, err)

	err = i9.Close()
	assert.Nil(t, err)

	err = i10.Close()
	assert.Nil(t, err)

	// fail: multiple request with different config of the same line
	i8, err = core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i8)
	i9, err = core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateActive))
	assert.Nil(t, err)
	require.NotNil(t, i9)
	i10, err = core.RequestItem(platform.Devpath(), io, core.AsInput(core.PullDisabled))
	assert.NotNil(t, err)
	require.Nil(t, i10)

	err = i8.Close()
	assert.Nil(t, err)

	err = i9.Close()
	assert.Nil(t, err)

	require.True(t, i9.Closed())
}

func TestItemClose(t *testing.T) {
	// success: single owner close
	i := newItem(t)
	assert.False(t, i.Closed())
	err := i.Close()
	assert.Nil(t, err)
	assert.True(t, i.Closed())

	// success: multi owner close
	i2 := newItem(t)
	assert.False(t, i2.Closed())
	i3 := newItem(t)
	assert.False(t, i3.Closed())
	assert.Equal(t, i2, i3)
	i4 := newItem(t)
	assert.False(t, i4.Closed())
	assert.Equal(t, i3, i4)

	err = i2.Close()
	assert.Nil(t, err)
	assert.False(t, i2.Closed())
	err = i3.Close()
	assert.Nil(t, err)
	assert.False(t, i3.Closed())
	err = i4.Close()
	assert.Nil(t, err)
	assert.True(t, i4.Closed())
}

func TestNewWatcher(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// fail: bad chip name
	w, err := core.NewWatcher(platform.Devpath()+"not", io)
	assert.NotNil(t, err)
	require.Nil(t, w)

	// fail: negative Item offset
	w, err = core.NewWatcher(platform.Devpath(), -1)
	assert.NotNil(t, err)
	require.Nil(t, w)

	// fail: out-of-range Item offset
	w, err = core.NewWatcher(platform.Devpath(), platform.Lines())
	assert.NotNil(t, err)
	require.Nil(t, w)

	// success: without options
	w2, err := core.NewWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, w2)
	err = w2.Close()
	assert.Nil(t, err)

	// success: with input(core.PullDisabled) option
	w3, err := core.NewWatcher(platform.Devpath(), io, core.AsInput(core.PullDisabled))
	assert.Nil(t, err)
	require.NotNil(t, w3)
	err = w3.Close()
	assert.Nil(t, err)

	// success: with input(core.PullUp) option
	w4, err := core.NewWatcher(platform.Devpath(), io, core.AsInput(core.PullUp))
	assert.Nil(t, err)
	require.NotNil(t, w4)
	err = w4.Close()
	assert.Nil(t, err)

	// success: with input(core.PullDown) option
	w5, err := core.NewWatcher(platform.Devpath(), io, core.AsInput(core.PullDown))
	assert.Nil(t, err)
	require.NotNil(t, w5)
	err = w5.Close()
	assert.Nil(t, err)

	// success: with output(high)
	w6, err := core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateActive))
	assert.Nil(t, err)
	require.NotNil(t, w6)
	err = w6.Close()
	assert.Nil(t, err)

	// success: with output(low)
	w7, err := core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, w7)
	err = w7.Close()
	assert.Nil(t, err)

	// success: multiple request with same config of the same line
	w8, err := core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, w8)
	w9, err := core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, w9)
	w10, err := core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, w10)

	err = w8.Close()
	assert.Nil(t, err)
	err = w9.Close()
	assert.Nil(t, err)
	err = w10.Close()
	assert.Nil(t, err)

	// fail: multiple request with different config of the same line
	w8, err = core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, w8)
	w9, err = core.NewWatcher(platform.Devpath(), io, core.AsOutput(core.StateActive))
	assert.Nil(t, err)
	require.NotNil(t, w9)
	w10, err = core.NewWatcher(platform.Devpath(), io, core.AsInput(core.PullDisabled))
	assert.NotNil(t, err)
	require.Nil(t, w10)

	err = w8.Close()
	assert.Nil(t, err)
	err = w9.Close()
	assert.Nil(t, err)

	require.True(t, w9.Closed())
}

func TestInputWatcher(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// fail: bad chip name
	w, err := core.NewInputWatcher(platform.Devpath()+"not", io)
	assert.NotNil(t, err)
	require.Nil(t, w)

	// fail: negative Item offset
	w, err = core.NewInputWatcher(platform.Devpath(), -1)
	assert.NotNil(t, err)
	require.Nil(t, w)

	// fail: out-of-range Item offset
	w, err = core.NewInputWatcher(platform.Devpath(), platform.Lines())
	assert.NotNil(t, err)
	require.Nil(t, w)

	// success: single valid InputWatcher
	w2, err := core.NewInputWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, w2)
	err = w2.Close()
	assert.Nil(t, err)

	// success: multiple valid InputWatcher
	w8, err := core.NewInputWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, w8)
	w9, err := core.NewInputWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, w9)
	w10, err := core.NewInputWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, w10)

	err = w8.Close()
	assert.Nil(t, err)
	err = w9.Close()
	assert.Nil(t, err)
	err = w10.Close()
	assert.Nil(t, err)

	require.True(t, w10.Closed())
}

func TestWatch(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]
	watcher, err := core.NewWatcher(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, watcher)

	i := newItem(t)

	err = i.SetState(core.StateActive)
	assert.Nil(t, err)
	le, ok := <-watcher.Watch()
	assert.True(t, ok)
	require.NotNil(t, le)
	assert.False(t, le.IsLineEvent)

	err = i.SetState(core.StateInactive)
	assert.Nil(t, err)
	le, ok = <-watcher.Watch()
	assert.True(t, ok)
	require.NotNil(t, le)
	assert.False(t, le.IsLineEvent)

	err = i.SetState(core.StateActive)
	assert.Nil(t, err)

	le, ok = <-watcher.Watch()
	assert.True(t, ok)
	require.NotNil(t, le)
	assert.False(t, le.IsLineEvent)

	err = i.SetState(core.StateActive)
	assert.Nil(t, err)

	select {
	case <-watcher.Watch():
		assert.Fail(t, "le should not receive anymore messages")

	case le, ok = <-watcher.Watch():
		assert.Fail(t, "channel should not be closed yet")

	default:
	}

	err = core.Close()
	assert.Nil(t, err)

	select {
	case _, ok = <-watcher.Watch():
		assert.False(t, ok)

	default:
		assert.Fail(t, "channel should be closed")
	}
}

func TestGetItem(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// fail: bad chip name
	i, err := core.GetItem(platform.Devpath()+"not", io)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// fail: negative Item offset
	i, err = core.GetItem(platform.Devpath(), -1)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// fail: out-of-range Item offset
	i, err = core.GetItem(platform.Devpath(), platform.Lines())
	assert.NotNil(t, err)
	require.Nil(t, i)

	// fail: get item that hasn't been registered yet
	i, err = core.GetItem(platform.Devpath(), io)
	assert.NotNil(t, err)
	require.Nil(t, i)

	// success: get an item that has been registered
	i, err = core.RequestItem(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, i)
	i2, err := core.GetItem(platform.Devpath(), io)
	assert.Nil(t, err)
	require.NotNil(t, i2)
	assert.Equal(t, i, i2)

	err = i.Close()
	assert.Nil(t, err)

	require.True(t, i.Closed())
}

func TestSetState(t *testing.T) {
	// Item offset
	io := platform.FloatingLines()[0]

	// fail: bad chip name
	err := core.SetState(platform.Devpath()+"not", io, core.StateActive)
	assert.NotNil(t, err)

	// fail: negative Item offset
	err = core.SetState(platform.Devpath(), -1, core.StateActive)
	assert.NotNil(t, err)

	// fail: out-of-range Item offset
	err = core.SetState(platform.Devpath(), platform.Lines(), core.StateActive)
	assert.NotNil(t, err)

	// fail: core.SetState an item that hasn't been registered yet
	err = core.SetState(platform.Devpath(), io, core.StateActive)
	assert.NotNil(t, err)

	// success: core.SetState an item that has been registered
	i, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i)
	for _, s := range []core.State{core.StateActive, core.StateInactive, core.StateActive, core.StateActive} {
		err = core.SetState(platform.Devpath(), io, s)
		assert.Nil(t, err)
		assert.Equal(t, s, i.State())
	}

	err = i.Close()
	assert.Nil(t, err)

	require.True(t, i.Closed())
}

func TestClose(t *testing.T) {
	// success: single owner close
	i := newItem(t)
	assert.False(t, i.Closed())

	err := i.Close()
	assert.Nil(t, err)
	assert.True(t, i.Closed())

	err = core.Close()
	assert.Nil(t, err)

	assert.True(t, i.Closed())

	// success: multi owner close
	i2 := newItem(t)
	assert.False(t, i2.Closed())
	i3 := newItem(t)
	assert.False(t, i3.Closed())
	assert.Equal(t, i2, i3)
	i4 := newItem(t)
	assert.False(t, i4.Closed())
	assert.Equal(t, i3, i4)

	err = core.Close()
	assert.Nil(t, err)
	assert.True(t, i2.Closed())
	assert.True(t, i3.Closed())
	assert.True(t, i4.Closed())
}

func TestSetLogger(t *testing.T) {
	err := core.SetLogger(nil)
	assert.NotNil(t, err)

	err = core.SetLogger(core.DummyLogger{})
	assert.Nil(t, err)
}

func newItem(t *testing.T) core.Item {
	// Item offset
	io := platform.FloatingLines()[0]

	i, err := core.RequestItem(platform.Devpath(), io, core.AsOutput(core.StateInactive))
	assert.Nil(t, err)
	require.NotNil(t, i)

	return i
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
