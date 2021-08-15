package security

import (
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

var security registry = registry{}

type Security struct {
	state     core.State
	sensors   *itemRegistry
	actuators *itemRegistry

	mu *sync.RWMutex
}

func Register(tag string, opts ...Option) (s *Security, err error) {
	options := &Options{}
	for _, opt := range opts {
		err = opt.applyOption(options)
		if err != nil {
			return
		}
	}
	s = &Security{
		state:     core.Inactive,
		sensors:   &itemRegistry{},
		actuators: &itemRegistry{},
		mu:        &sync.RWMutex{},
	}
	for chip, opt := range options.control {
		err = s.AddSensor(chip, tag, opt.sensors)
		if err != nil {
			return
		}
		err = s.AddActuator(chip, tag, opt.actuators)
		if err != nil {
			return
		}
	}
	return
}

func (s *Security) SetState(state core.State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
	s.actuators.ForEach(func(i *core.Item) {
		i.SetState(state)
	})
}

func (s *Security) AddSensor(gpioName string, tag string, offsets []int) (err error) {
	for _, offset := range offsets {
		i, err := core.RegisterItem(gpioName, offset, core.AsInput(core.PullDown), core.WithState(s.state))
		if err != nil {
			return err
		}
		i.AddEventListener(func(event *core.ItemEvent) {
			if event.Item.State() == core.Active {
				s.SetState(core.Active)
			}
		})
		s.sensors.Add(gpioName, offset, i)
	}
	return
}

func (s *Security) AddActuator(gpioName string, tag string, offsets []int) (err error) {
	for _, offset := range offsets {
		i, err := core.RegisterItem(gpioName, offset, core.AsOutput(), core.WithState(core.Inactive))
		if err != nil {
			return err
		}
		s.actuators.Add(gpioName, offset, i)
	}
	return
}
