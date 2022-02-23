package general

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"go.uber.org/multierr"
)

type Sensor = core.Watcher

type Actuator interface {
	Close() error
	SetState(core.State) error
	State() core.State
}

type GeneralI interface {
	Close() error
	Register(tag string, opts ...Option) (g *General, err error)
	State() core.State
}
type General struct {
	state     core.State
	sensors   []Sensor
	actuators []Actuator
	active    bool
	kind      string
	strategy  string

	*sync.RWMutex
}

func Register(tag string, opts ...Option) (g *General, err error) {
	options := &Options{
		control: map[string]Control{},
	}
	for _, opt := range opts {
		err = opt.applyOption(options)
		if err != nil {
			return
		}
	}

	g = &General{
		state:     core.StateInactive,
		sensors:   []Sensor{},
		actuators: []Actuator{},
		active:    true,
		kind:      options.kind,
		strategy:  options.strategy,
		RWMutex:   &sync.RWMutex{},
	}

	for chip, opt := range options.control {
		err = g.AddSensor(chip, opt.sensors...)
		if err != nil {
			return
		}
		err = g.AddActuator(chip, opt.actuators...)
		if err != nil {
			return
		}
	}

	initalState := core.StateInactive
	if options.kind == RSync {
		initalState = core.StateActive
	}
	g.SetState(initalState)

	err = g.handle()
	if err != nil {
		return nil, err
	}

	return
}

func (g *General) handle() error {
	var handler func(core.EventChannel)
	switch g.kind {
	case Alarm:
		handler = g.AlarmHandler
	case Sync:
		switch g.strategy {
		case AllIn:
			handler = g.SyncHandlerAllIn
		case OneIn:
			handler = g.SyncHandlerOneIn
		}
	case RSync:
		switch g.strategy {
		case AllIn:
			handler = g.RSyncHandlerAllIn
		case OneIn:
			handler = g.RSyncHandlerOneIn
		}
	default:
		return fmt.Errorf("mode should be set")
	}

	for _, sensor := range g.sensors {
		go handler(sensor.Watch())
	}
	return nil
}

func (g *General) State() core.State {
	g.Lock()
	defer g.Unlock()
	return g.state
}

func (g *General) SetState(state core.State) (err error) {
	g.Lock()
	if state == g.state {
		g.Unlock()
		return
	}
	g.state = state
	actuators := g.actuators
	g.Unlock()

	for _, actuator := range actuators {
		err = multierr.Append(err, actuator.SetState(state))
	}
	return
}

func (g *General) AddSensor(gpioName string, offsets ...int) (err error) {
	for _, offset := range offsets {
		var sensor Sensor
		sensor, err = core.NewWatcher(gpioName, offset, core.AsInput(core.PullDown), core.WithState(g.state))
		if err != nil {
			return
		}
		g.sensors = append(g.sensors, sensor)
	}
	return
}

func (g *General) AddActuator(gpioName string, offsets ...int) (err error) {
	for _, offset := range offsets {
		var actuator Actuator
		actuator, err = core.RequestItem(gpioName, offset, core.AsOutput(core.StateInactive))
		if err != nil {
			return err
		}
		g.actuators = append(g.actuators, actuator)
	}
	return
}

func (g *General) SetActive(state bool) {
	if !state {
		g.SetState(core.StateInactive)
	}
	g.active = state
}

func (g *General) AlarmHandler(ch core.EventChannel) {
	for ie := range ch {
		if ie.Info.State == core.StateActive {
			g.SetState(core.StateActive)
		}
	}
}

func (g *General) SyncHandlerAllIn(ch core.EventChannel) {
	for range ch {
		g.Lock()
		sensors := g.sensors
		g.Unlock()
		allIn := true
		for _, sensor := range sensors {
			if sensor.State() == core.StateInactive {
				allIn = false
			}
		}
		if allIn {
			g.SetState(core.StateActive)
		} else {
			g.SetState(core.StateInactive)
		}
	}
}

func (g *General) SyncHandlerOneIn(ch core.EventChannel) {
	for range ch {
		g.Lock()
		sensors := g.sensors
		g.Unlock()
		oneIn := false
		for _, sensor := range sensors {
			if sensor.State() == core.StateActive {
				oneIn = true
			}
		}
		if oneIn {
			g.SetState(core.StateActive)
		} else {
			g.SetState(core.StateInactive)
		}
	}
}

func (g *General) RSyncHandlerAllIn(ch core.EventChannel) {
	for range ch {
		g.Lock()
		sensors := g.sensors
		g.Unlock()
		allIn := true
		for _, sensor := range sensors {
			if sensor.State() == core.StateActive {
				allIn = false
			}
		}
		if allIn {
			g.SetState(core.StateActive)
		} else {
			g.SetState(core.StateInactive)
		}
	}
}

func (g *General) RSyncHandlerOneIn(ch core.EventChannel) {
	for range ch {
		g.Lock()
		sensors := g.sensors
		g.Unlock()
		oneIn := false
		for _, sensor := range sensors {
			if sensor.State() == core.StateInactive {
				oneIn = true
			}
		}
		if oneIn {
			g.SetState(core.StateActive)
		} else {
			g.SetState(core.StateInactive)
		}
	}
}
