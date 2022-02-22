package core

import (
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
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
		c.RLock()
		c.Chip = chip
		c.RUnlock()
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
	c.RLock()
	tgc := c.tgc
	c.RUnlock()
	return tgc.State()
}

func (c *chip) GetItem(offset int) (Item, error) {
	return c.getItem(offset)
}

func (c *chip) getItem(offset int) (*item, error) {
	c.RLock()
	defer c.RUnlock()
	return c.items.Get(offset)
}

func (c *chip) Name() string {
	c.RLock()
	defer c.RUnlock()
	return c.name
}

func (c *chip) Close() error {
	c.RLock()
	tgc := c.tgc
	c.RUnlock()

	tgc.Delete()

	return nil
}

func (c *chip) cleanup() (err error) {
	c.RLock()
	chip := c.Chip
	c.RUnlock()
	chips.Delete(chip.Name)
	err = chip.Close()
	if err != nil {
		logger.Errorf("something went wrong while closing chip %s: %v", chip.Name, err)
	}

	logger.Infof("%s is successfuly cleaned up", chip.Name)

	return
}
