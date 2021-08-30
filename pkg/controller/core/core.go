package core

import (
	"fmt"

	"go.uber.org/multierr"
)

func GetChip(chipName string) (c *Chip, err error) {
	return chips.Get(chipName)
}

func GetItem(chipName string, offset int) (i *Item, err error) {
	c, err := GetChip(chipName)
	if err != nil {
		return nil, err
	}
	return c.GetItem(offset)
}

func RegisterItem(chip string, offset int, opts ...ItemOption) (item *Item, err error) {
	// get the chip
	c, err := GetChip(chip)
	if err != nil {
		return nil, fmt.Errorf("there is no registered chip named %s", chip)
	}

	return c.RegisterItem(offset, opts...)
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

func AddEventListener(chipName string, offset int, fns ...EventHandler) (err error) {
	i, err := GetItem(chipName, offset)
	if err != nil {
		return
	}
	return i.AddEventListener(fns...)
}

func Close() (err error) {
	chips.ForEach(func(chipName string, chip *Chip) {
		err = multierr.Append(err, chip.Close())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("gpio core is successfully cleanedup")
	}
	return
}
