package core

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

var chips = newChipRegistry()

type ChipInfo struct {
}

type Chip interface {
	Closer
	RequestItem(offset int, opts ...ItemOption) (Item, error)
	Used() bool
	Info() (ChipInfo, error)
	GetItem(offset int) (Item, error)
}

type chip struct {
	*gpiod.Chip
	items *itemRegistry
	tgc   *tgc.Tgc
	name  string
	*sync.RWMutex
}

func RequestChip(name string) (Chip, error) {
	c, err := getChip(name)

	// if chip doesn't exits, create it
	if _, ok := err.(ChipNotFoundError); ok {
		c = &chip{
			Chip:    nil,
			items:   newItemRegistry(),
			tgc:     nil,
			name:    name,
			RWMutex: &sync.RWMutex{},
		}

		var t *tgc.Tgc
		t, err = tgc.New(c.tgcHandler)
		if err != nil {
			return nil, err
		}
		c.tgc = t

		err = chips.Add(name, c)
		if err != nil {
			return nil, err
		}

		logger.Infof("chip %s registerd successfully by %s", name)
	}

	// if chip exist just add new owner to it
	c.Lock()
	c.tgc.Add()
	c.Unlock()

	return c, err
}

func (c *chip) tgcHandler(b bool) {
	if b {
		chip, err := gpiod.NewChip(c.name)
		if err != nil {
			return
		}
		c.Lock()
		c.Chip = chip
		c.Unlock()
		chips.Add(c.name, c)
	} else {
		chips.Delete(c.name)
		c.cleanup()
	}
}

func (c *chip) RequestItem(offset int, opts ...ItemOption) (Item, error) {
	c.Lock()
	itemReg := c.items
	c.Unlock()

	// apply options
	options := &ItemOptions{}
	for _, io := range opts {
		err := io.applyItemOption(options)
		if err != nil {
			return nil, err
		}
	}

	i, err := itemReg.Get(offset)

	if _, ok := err.(ItemNotFound); ok {
		// item doesnt exist in registry so we'll create it
		i := &item{
			Line:    nil,
			RWMutex: &sync.RWMutex{},
			chip:    c,
			state:   options.state,
			offset:  offset,
			events:  newEventRegistry(),
			options: options,
			tgc:     nil,
		}

		t, err := tgc.New(i.tgcHandler)
		if err != nil {
			return nil, err
		}
		i.tgc = t

		err = itemReg.Add(offset, i)
		if err != nil {
			return nil, err
		}

		logger.Infof("item registerd on line %o as %s", offset, options.mode)
	} else {
		// already exits, check if its of the same line direction
		i.Lock()
		info, err := i.Line.Info()
		i.Unlock()
		if err != nil {
			return nil, err
		}
		if info.Config.Direction != gpiod.LineDirection(options.mode) {
			return nil, fmt.Errorf("this item is already registered as %s", Mode(info.Config.Direction))
		}
	}

	i.Lock()
	i.tgc.Add()
	i.Unlock()

	return i, nil
}

func (c *chip) Info() (ChipInfo, error) {
	return ChipInfo{}, nil
}

func (c *chip) Used() bool {
	c.Lock()
	defer c.Unlock()
	return c.tgc.State()
}

func (c *chip) GetItem(offset int) (Item, error) {
	return c.getItem(offset)
}

func (c *chip) getItem(offset int) (i *item, err error) {
	c.Lock()
	defer c.Unlock()
	return c.items.Get(offset)
}

func (c *chip) Close() error {
	c.Lock()
	defer c.Unlock()
	c.tgc.Delete()
	return nil
}

func (c *chip) cleanup() (err error) {
	c.Lock()
	ir := c.items
	multierr.Append(err, c.Chip.Close())
	chipName := c.Chip.Name
	c.Unlock()
	if err != nil {
		logger.Errorf(err.Error())
	}
	ir.ForEach(func(offset int, i *item) {
		err = multierr.Append(err, i.Close())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("%s is successfuly cleaned up", chipName)
	}
	return
}
