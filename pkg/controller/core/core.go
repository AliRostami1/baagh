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
var chips chipRegistry

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
		chip:  c,
		items: &itemRegistry{},
	}
	err = chips.append(options.name, chip)
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
		line:   nil,
		state:  options.state,
		events: &EventRegistry{},
		mu:     sync.RWMutex{},
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
	case Output:
		var l *gpiod.Line
		l, err = c.chip.RequestLine(offset, gpiod.AsOutput(int(options.state)))
		if err != nil {
			return nil, err
		}
		item.line = l
	default:
		err = fmt.Errorf("you have to set the mode")
	}
	return
}

func SetState(chipName string, offset int, state State) (err error) {
	c, err := chips.get(chipName)
	if err != nil {
		return
	}
	l, err := c.items.get(offset)
	if err != nil {
		return
	}
	return l.SetState(state)
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

	mu sync.RWMutex
}

func (i *Item) SetState(state State) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
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
