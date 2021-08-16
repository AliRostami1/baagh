package general

import (
	"log"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

var general registry = registry{
	registry: map[string]*General{},
	RWMutex:  &sync.RWMutex{},
}

type General struct {
	state     core.State
	sensors   *itemRegistry
	actuators *itemRegistry
	kind      string
	strategy  string

	mu *sync.RWMutex
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
		state: core.Inactive,
		sensors: &itemRegistry{
			registry: map[string]map[int]*core.Item{},
			RWMutex:  &sync.RWMutex{},
		},
		actuators: &itemRegistry{
			registry: map[string]map[int]*core.Item{},
			RWMutex:  &sync.RWMutex{},
		},
		kind:     options.kind,
		strategy: options.strategy,
		mu:       &sync.RWMutex{},
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
	return
}

func (g *General) setState(state core.State) {
	g.mu.Lock()
	g.state = state
	actuators := g.actuators
	g.mu.Unlock()
	actuators.ForEach(func(i *core.Item) {
		i.SetState(state)
	})
}

func (g *General) AddSensor(gpioName string, tag string, offsets []int) (err error) {
	log.Printf("*g = %#v", *g)
	for _, offset := range offsets {
		i, err := core.RegisterItem(gpioName, offset, core.AsInput(core.PullDown), core.WithState(g.state))
		if err != nil {
			return err
		}
		var handler core.EventHandler
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
		}
		if handler == nil {
			panic("mode should be set")
		}

		i.AddEventListener(handler)
		g.sensors.Add(gpioName, offset, i)
	}
	return
}

func (g *General) AddActuator(gpioName string, tag string, offsets []int) (err error) {
	for _, offset := range offsets {
		i, err := core.RegisterItem(gpioName, offset, core.AsOutput(), core.WithState(core.Inactive))
		if err != nil {
			return err
		}
		g.actuators.Add(gpioName, offset, i)
	}
	return
}

func (g *General) TurnOff() {
	g.setState(core.Inactive)
}

func (g *General) TurnOn() {
	g.setState(core.Active)
}

func (g *General) AlarmHandler(event *core.ItemEvent) {
	if event.Item.State() == core.Active {
		g.setState(core.Active)
	}
}

func (g *General) SyncHandlerAllIn(event *core.ItemEvent) {
	g.mu.Lock()
	sensors := g.sensors
	g.mu.Unlock()
	allIn := true
	sensors.ForEach(func(i *core.Item) {
		if i.State() == core.Inactive {
			allIn = false
		}
	})
	if allIn {
		g.TurnOn()
	} else {
		g.TurnOff()
	}
}

func (g *General) SyncHandlerOneIn(event *core.ItemEvent) {
	g.mu.Lock()
	sensors := g.sensors
	g.mu.Unlock()
	oneIn := false
	sensors.ForEach(func(i *core.Item) {
		if i.State() == core.Active {
			oneIn = true
		}
	})
	if oneIn {
		g.TurnOn()
	} else {
		g.TurnOff()
	}
}

func (g *General) RSyncHandlerAllIn(event *core.ItemEvent) {
	g.mu.Lock()
	sensors := g.sensors
	g.mu.Unlock()
	allIn := true
	sensors.ForEach(func(i *core.Item) {
		if i.State() == core.Active {
			allIn = false
		}
	})
	if allIn {
		g.TurnOn()
	} else {
		g.TurnOff()
	}
}

func (g *General) RSyncHandlerOneIn(event *core.ItemEvent) {
	g.mu.Lock()
	sensors := g.sensors
	g.mu.Unlock()
	oneIn := false
	sensors.ForEach(func(i *core.Item) {
		if i.State() == core.Inactive {
			oneIn = true
		}
	})
	if oneIn {
		g.TurnOn()
	} else {
		g.TurnOff()
	}
}
