package core

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

var chips = newChipRegistry()

type Chip struct {
	chip    *gpiod.Chip
	items   *itemRegistry
	tgc     *tgc.Tgc
	options *ChipOptions

	mu *sync.RWMutex
}

func RegisterChip(name string, opts ...ChipOption) (chip *Chip, err error) {
	options := &ChipOptions{
		name:     name,
		consumer: defaultConsumer(),
	}
	for _, co := range opts {
		err = co.applyChipOption(options)
		if err != nil {
			return
		}
	}

	chip, err = GetChip(name)

	if _, ok := err.(ChipNotFoundError); ok {
		chip = &Chip{
			chip:    nil,
			items:   newItemRegistry(),
			tgc:     tgc.New(chip.tgcHandler),
			options: options,
			mu:      &sync.RWMutex{},
		}
		logger.Infof("chip %s registerd successfully by %s", name, options.consumer)
		return
	}

	if err != nil {
		return
	}

	chip.mu.Lock()
	chip.tgc.Add()
	chip.mu.Unlock()

	return
}

func (c *Chip) tgcHandler(b bool) {
	if b {
		c.createChip()
	} else {
		c.removeChip()
	}
}

func (c *Chip) removeChip() (err error) {
	err = c.Cleanup()
	c.chip = nil
	chips.Delete(c.options.name)
	return
}

func (c *Chip) createChip() (err error) {
	chip, err := gpiod.NewChip(c.options.name, gpiod.WithConsumer(c.options.consumer))
	if err != nil {
		return
	}
	c.mu.Lock()
	c.chip = chip
	c.mu.Unlock()
	err = chips.Add(c.options.name, c)
	if err != nil {
		return err
	}
	return
}

func (c *Chip) RegisterItem(offset int, opts ...ItemOption) (item *Item, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// apply options
	options := &ItemOptions{}
	for _, io := range opts {
		err = io.applyItemOption(options)
		if err != nil {
			return nil, err
		}
	}

	item, err = c.items.Get(offset)

	if _, ok := err.(ItemNotFound); ok {
		// item doesnt exist in registry so we'll create it
		item = newItem(options.state)
		item.AddEventListener(func(event *ItemEvent) {
			events.CallAll(event)
		})

		switch options.mode {
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
			return nil, fmt.Errorf("you have to set the mode")
		}

		err = c.items.Add(offset, item)
		if err != nil {
			return nil, err
		}
		logger.Infof("item registerd on line %o as %s", offset, options.mode)
		return
	}

	// errors other than ItemNotFound shoudd just be returned
	if err != nil {
		return nil, err
	}

	// already exits, check if its of the same line direction
	info, err := item.line.Info()
	if err != nil {
		return nil, err
	}
	if info.Config.Direction != gpiod.LineDirection(options.mode) {
		return nil, fmt.Errorf("this item is already registered as %s", Mode(info.Config.Direction))
	}
	item.incrOwner()
	return item, nil
}

func (c *Chip) GetItem(offset int) (i *Item, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.items.Get(offset)
}

func (c *Chip) Cleanup() (err error) {
	c.mu.Lock()
	ir := c.items
	multierr.Append(err, c.chip.Close())
	chipName := c.chip.Name
	c.mu.Unlock()
	if err != nil {
		logger.Errorf(err.Error())
	}
	ir.ForEach(func(offset int, item *Item) {
		err = multierr.Append(err, item.Cleanup())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("%s is successfuly cleaned up", chipName)
	}
	return
}
