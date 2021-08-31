package core

import (
	"fmt"

	"go.uber.org/multierr"
)

type Closer interface {
	Close() error
}

func GetChip(chipName string) (c *chip, err error) {
	return chips.Get(chipName)
}

func GetItem(chipName string, offset int) (i *item, err error) {
	c, err := GetChip(chipName)
	if err != nil {
		return nil, err
	}
	return c.GetItem(offset)
}

func RequestItem(chip string, offset int, opts ...ItemOption) (item *item, err error) {
	// get the chip
	c, err := GetChip(chip)
	if err != nil {
		return nil, fmt.Errorf("there is no registered chip named %s", chip)
	}

	return c.RequestItem(offset, opts...)
}

func SetState(chipName string, offset int, state State) (err error) {
	i, err := GetItem(chipName, offset)
	if err != nil {
		return
	}
	err = i.SetState(state)
	if err != nil {
		return
	}
	return
}

func NewWatcher(chipName string, offset int) (Watcher, error) {
	i, err := GetItem(chipName, offset)
	if err != nil {
		return nil, err
	}
	return i.NewWatcher(), nil
}

func Close() (err error) {
	chips.ForEach(func(chipName string, chip *chip) {
		err = multierr.Append(err, chip.cleanup())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("gpio core is successfully cleanedup")
	}
	return
}
