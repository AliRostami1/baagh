package general

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

var general = registry{
	registry: map[string]*General{},
	RWMutex:  &sync.RWMutex{},
}

type GeneralI interface {
	core.Closer
	Register(tag string, opts ...Option) (g *General, err error)
	State() core.State
}
type General struct {
	state     core.State
	sensors   *itemRegistry
	actuators *itemRegistry
	active    bool
	kind      string
	strategy  string
	watchers  []core.Watcher

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
		sensors:   &itemRegistry{registry: map[string]map[int]core.Item{}, RWMutex: &sync.RWMutex{}},
		actuators: &itemRegistry{registry: map[string]map[int]core.Item{}, RWMutex: &sync.RWMutex{}},
		active:    true,
		kind:      options.kind,
		strategy:  options.strategy,
		watchers:  []core.Watcher{},
		RWMutex:   &sync.RWMutex{},
	}

	for chip, opt := range options.control {
		err = g.AddSensor(chip, tag, opt.sensors)
		if err != nil {
			return
		}
		err = g.AddActuator(chip, tag, opt.actuators)
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

	for _, watcher := range g.watchers {
		go handler(watcher.Watch())
	}
	return nil
}

func (g *General) State() core.State {
	g.Lock()
	defer g.Unlock()
	return g.state
}

func (g *General) SetState(state core.State) {
	g.Lock()
	if state == g.state {
		g.Unlock()
		return
	}
	g.state = state
	actuators := g.actuators
	g.Unlock()

	actuators.ForEach(func(i core.Item) {
		i.SetState(state)
	})
}

func (g *General) AddSensor(gpioName string, tag string, offsets []int) (err error) {
	for _, offset := range offsets {
		watcher, err := core.NewWatcher(gpioName, offset, core.AsInput(core.PullDown), core.WithState(g.state))
		if err != nil {
			return err
		}

		g.watchers = append(g.watchers, watcher)
		g.sensors.Add(gpioName, offset, i)
	}
	return
}

func (g *General) AddActuator(gpioName string, tag string, offsets []int) (err error) {
	for _, offset := range offsets {
		i, err := core.RequestItem(gpioName, offset, core.AsOutput(core.StateInactive))
		if err != nil {
			return err
		}
		g.actuators.Add(gpioName, offset, i)
	}
	return
}

func (g *General) SetStateActive(a bool) {
	g.active = a
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
		sensors.ForEach(func(i core.Item) {
			if i.State() == core.StateInactive {
				allIn = false
			}
		})
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
		sensors.ForEach(func(i core.Item) {
			if i.State() == core.StateActive {
				oneIn = true
			}
		})
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
		sensors.ForEach(func(i core.Item) {
			if i.State() == core.StateActive {
				allIn = false
			}
		})
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
		sensors.ForEach(func(i core.Item) {
			if i.State() == core.StateInactive {
				oneIn = true
			}
		})
		if oneIn {
			g.SetState(core.StateActive)
		} else {
			g.SetState(core.StateInactive)
		}
	}
}
