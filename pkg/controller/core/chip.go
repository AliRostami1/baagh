package core

import (
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

var chips = newChipRegistry()

// TODO: implement this
type ChipInfo struct {
}

type chip struct {
	*gpiod.Chip
	items *itemRegistry
	tgc   *tgc.Tgc
	name  string
	*sync.RWMutex
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

func (c *chip) Info() (ChipInfo, error) {
	return ChipInfo{}, nil
}

func (c *chip) Used() bool {
	c.Lock()
	tgc := c.tgc
	c.Unlock()
	return tgc.State()
}

func (c *chip) GetItem(offset int) (Item, error) {
	return c.getItem(offset)
}

func (c *chip) getItem(offset int) (i *item, err error) {
	c.Lock()
	defer c.Unlock()
	return c.items.Get(offset)
}

func (c *chip) Name() string {
	c.Lock()
	defer c.Unlock()
	return c.name
}

func (c *chip) Close() error {
	c.Lock()
	tgc := c.tgc
	c.Unlock()
	tgc.Delete()
	return nil
}

func (c *chip) cleanup() (err error) {
	c.Lock()
	ir := c.items
	multierr.Append(err, c.Chip.Close())
	chipName := c.Chip.Name
	c.Unlock()
	ir.ForEach(func(offset int, i *item) {
		err = multierr.Append(err, i.cleanup())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("%s is successfuly cleaned up", chipName)
	}
	return
}
