package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/log"

	"github.com/warthog618/gpiod"
)

// key is chip name
var chips chipRegistry

var logger log.Logger

var ctx context.Context = context.Background()

type Chip struct {
	chip  *gpiod.Chip
	items *itemRegistry
}

func SetLogger(l log.Logger) {
	logger = l
}

func RegisterChip(ctx context.Context, opts ...ChipOption) (err error) {
	options := &ChipOptions{}
	for _, co := range opts {
		err = co.applyChipOption(options)
		if err != nil {
			return err
		}
	}
	c, err := gpiod.NewChip(options.name, gpiod.WithConsumer(options.consumer))
	if err != nil {
		return err
	}
	return chips.append(options.name, &Chip{
		chip:  c,
		items: &itemRegistry{},
	})
}

type Item struct {
	line   *gpiod.Line
	state  State
	events *EventRegistry

	mu sync.RWMutex
}

func (i *Item) SetState(state State) (err error) {
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

func RegisterItem(chip string, pin int, opts ...ItemOption) (err error) {
	// apply options
	options := &ItemOptions{}
	for _, io := range opts {
		err = io.applyItemOption(options)
		if err != nil {
			return err
		}
	}

	// get the chip
	c, err := chips.get(chip)
	if err != nil {
		return fmt.Errorf("there is no registered chip named %s", chip)
	}

	item := &Item{
		line:   nil,
		state:  options.state,
		events: &EventRegistry{},
		mu:     sync.RWMutex{},
	}

	item.mu.Lock()
	defer item.mu.Unlock()

	if options.io.mode == Input {
		handler := func(evt gpiod.LineEvent) {
			if evt.Type == gpiod.LineEventRisingEdge {
				item.SetState(Active)
				return
			}
			if evt.Type == gpiod.LineEventFallingEdge {
				item.SetState(Inactive)
				return
			}
		}
		var l *gpiod.Line
		l, err = c.chip.RequestLine(pin, gpiod.AsInput, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
		if err != nil {
			return
		}
		item.line = l
		return
	}

	if options.io.mode == Output {

	}
	return
}

func SetState(chipName string, pin int, state State) (err error) {
	_, err = chips.get(chipName)
	if err != nil {
		return
	}
	// ...
	return
}
