package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/log"

	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

// key is chip name
var chips chipRegistry = chipRegistry{
	registry: map[string]*Chip{},
	RWMutex:  &sync.RWMutex{},
}

var logger log.Logger = log.DummyLogger{}

var ctx context.Context = context.Background()

func SetLogger(l log.Logger) {
	logger = l
}

func RegisterChip(ctx context.Context, opts ...ChipOption) (chip *Chip, err error) {
	options := &ChipOptions{}
	for _, co := range opts {
		err = co.applyChipOption(options)
		if err != nil {
			return
		}
	}
	c, err := gpiod.NewChip(options.name, gpiod.WithConsumer(options.consumer))
	if err != nil {
		return
	}
	chip = &Chip{
		chip: c,
		items: &itemRegistry{
			registry: map[int]*Item{},
			RWMutex:  &sync.RWMutex{},
		},
	}
	err = chips.append(options.name, chip)
	logger.Infof("chip %s registerd successfully by %s", options.name, options.consumer)
	return
}

func RegisterItem(chip string, offset int, opts ...ItemOption) (item *Item, err error) {
	// apply options
	options := &ItemOptions{}
	for _, io := range opts {
		err = io.applyItemOption(options)
		if err != nil {
			return nil, err
		}
	}

	// get the chip
	c, err := chips.get(chip)
	if err != nil {
		return nil, fmt.Errorf("there is no registered chip named %s", chip)
	}

	item = &Item{
		line:  nil,
		state: options.state,
		events: &EventRegistry{
			events:  []EventHandler{},
			RWMutex: &sync.RWMutex{},
		},
		mu: &sync.RWMutex{},
	}

	item.mu.Lock()
	defer item.mu.Unlock()

	switch options.io.mode {
	case Input:
		handler := func(evt gpiod.LineEvent) {
			switch evt.Type {
			case gpiod.LineEventRisingEdge:
				item.SetState(Active)
			case gpiod.LineEventFallingEdge:
				item.SetState(Inactive)
			}
		}
		var l *gpiod.Line
		l, err = c.chip.RequestLine(offset, gpiod.AsInput, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
		if err != nil {
			return nil, err
		}
		item.line = l
		logger.Infof("item registerd on line %o as %s", offset, options.io.mode)
	case Output:
		var l *gpiod.Line
		l, err = c.chip.RequestLine(offset, gpiod.AsOutput(int(options.state)))
		if err != nil {
			return nil, err
		}
		item.line = l
		logger.Infof("item registerd on line %o as %s", offset, options.io.mode)
	default:
		err = fmt.Errorf("you have to set the mode")
		logger.Errorf(err.Error())
	}
	return
}

func SetState(chipName string, offset int, state State) (err error) {
	c, err := chips.get(chipName)
	if err != nil {
		return
	}
	i, err := c.items.get(offset)
	if err != nil {
		return
	}
	err = i.SetState(state)
	if err != nil {
		return
	}
	return
}

func AddEventListener(chipName string, offset int, fns ...EventHandler) (err error) {
	c, err := chips.get(chipName)
	if err != nil {
		return
	}
	i, err := c.items.get(offset)
	if err != nil {
		return
	}
	return i.AddEventListener(fns...)
}

func Cleanup() (err error) {
	chips.forEach(func(chipName string, chip *Chip) {
		err = multierr.Append(err, chip.Cleanup())
	})
	return
}

type Chip struct {
	chip  *gpiod.Chip
	items *itemRegistry
}

func (c *Chip) Cleanup() (err error) {
	c.items.forEach(func(offset int, item *Item) {
		err = multierr.Append(err, item.Cleanup())
	})
	multierr.Append(err, c.chip.Close())
	return
}

type Item struct {
	line   *gpiod.Line
	state  State
	events *EventRegistry

	mu *sync.RWMutex
}

func (i *Item) SetState(state State) (err error) {
	i.mu.Lock()
	info, err := i.line.Info()
	if err != nil {
		return
	}
	if info.Config.Direction == gpiod.LineDirectionOutput {
		err = i.line.SetValue(int(state))
		if err != nil {
			return
		}
	}
	i.state = state
	events := i.events
	i.mu.Unlock()

	events.callAll(&ItemEvent{
		Item: i,
	})
	logger.Debugf("state changed to %s on line %o of chip %s", state, i.line.Offset(), i.line.Chip())
	return
}

func (i *Item) State() State {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.state
}

func (i *Item) AddEventListener(fns ...EventHandler) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	err = i.events.addEventListener(fns...)
	return
}

func (i *Item) Cleanup() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.line.Close()
}
