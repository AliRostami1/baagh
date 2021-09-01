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
		state:     core.Inactive,
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

	initalState := core.Inactive
	if options.kind == RSync {
		initalState = core.Active
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
		i, err := core.RequestItem(gpioName, offset, core.AsInput(core.PullDown), core.WithState(g.state))
		if err != nil {
			return err
		}

		watcher, err := i.NewWatcher()
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
		i, err := core.RequestItem(gpioName, offset, core.AsOutput(core.Inactive))
		if err != nil {
			return err
		}
		g.actuators.Add(gpioName, offset, i)
	}
	return
}

func (g *General) SetActive(a bool) {
	g.active = a
}

func (g *General) AlarmHandler(ch core.EventChannel) {
	for ie := range ch {
		if ie.Info.State == core.Active {
			g.SetState(core.Active)
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
			if i.State() == core.Inactive {
				allIn = false
			}
		})
		if allIn {
			g.SetState(core.Active)
		} else {
			g.SetState(core.Inactive)
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
			if i.State() == core.Active {
				oneIn = true
			}
		})
		if oneIn {
			g.SetState(core.Active)
		} else {
			g.SetState(core.Inactive)
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
			if i.State() == core.Active {
				allIn = false
			}
		})
		if allIn {
			g.SetState(core.Active)
		} else {
			g.SetState(core.Inactive)
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
			if i.State() == core.Inactive {
				oneIn = true
			}
		})
		if oneIn {
			g.SetState(core.Active)
		} else {
			g.SetState(core.Inactive)
		}
	}
}
