package core

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

type Closer interface {
	Close() error
}

func getItem(chip string, offset int) (i *item, err error) {
	c, err := chips.Get(chip)
	if err != nil {
		return nil, err
	}
	return c.getItem(offset)
}

func GetItem(chip string, offset int) (i Item, err error) {
	return getItem(chip, offset)
}

func requestChip(name string) (*chip, error) {
	c, err := chips.Get(name)

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

		logger.Infof("chip %[1]s registerd successfully", name)
	} else {
		logger.Debugf("chip %[1]s got a new owner", name)
	}

	// either if it just got created or it was there all along, increment it's tgc
	c.RLock()
	tgc := c.tgc
	c.RUnlock()
	tgc.Add()

	return c, err
}

func requestItem(chip string, offset int, opts ...ItemOption) (*item, error) {
	c, err := requestChip(chip)
	if err != nil {
		return nil, err
	}

	c.RLock()
	itemReg, chipName := c.items, c.name
	c.RUnlock()

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
		i = &item{
			Line:    nil,
			RWMutex: &sync.RWMutex{},
			chip:    c,
			state:   options.state,
			offset:  offset,
			events:  newEventRegistry(),
			options: options,
			tgc:     nil,
			closed:  false,
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

		logger.Infof("item registerd on line %d of chip %s as %s", offset, chipName, options.mode)
	} else {
		// already exits, check if its of the same line direction
		i.RLock()
		info, err := i.Line.Info()
		i.RUnlock()
		if err != nil {
			return nil, err
		}
		if info.Config.Direction != gpiod.LineDirection(options.mode) {
			return nil, fmt.Errorf("this item is already registered as %s, you can't register it as %s", Mode(info.Config.Direction), options.mode)
		}
		logger.Debugf("item registerd on line %d of chip %s as %s got a new owner", offset, chipName, options.mode)
	}

	// either if it just got created or it was there all along, increment it's tgc
	i.RLock()
	tgc := i.tgc
	i.RUnlock()
	tgc.Add()

	return i, nil
}

func RequestItem(chip string, offset int, opts ...ItemOption) (Item, error) {
	return requestItem(chip, offset, opts...)
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

func NewWatcher(chipName string, offset int, opts ...ItemOption) (Watcher, error) {
	i, err := requestItem(chipName, offset, opts...)
	if err != nil {
		return nil, err
	}

	w := &watcher{
		item:         i,
		eventChannel: make(chan *ItemEvent),
	}

	i.RLock()
	ev := i.events
	i.RUnlock()
	ev.Add(w.eventChannel)

	return w, nil
}

func NewInputWatcher(chipName string, offset int) (Watcher, error) {
	i, err := requestItem(chipName, offset, AsInput(PullDown))
	if err != nil {
		return nil, err
	}

	i.RLock()
	ev := i.events
	i.RUnlock()

	w := &watcher{
		item:         i,
		eventChannel: make(chan *ItemEvent),
	}

	ev.Add(w.eventChannel)

	return w, nil
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

func Chips() []string {
	return gpiod.Chips()
}
