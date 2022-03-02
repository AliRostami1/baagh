package core

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

var reg = registry{
	chips:   map[string]map[int]*item{},
	lines:   map[string]int{},
	RWMutex: &sync.RWMutex{},
}

func isChip(chip string) bool {
	err := gpiod.IsChip(chip)
	return err == nil
}

func isOffset(chip string, offset int) bool {
	var lines int
	if l, ok := reg.lines[chip]; !ok {
		c, _ := gpiod.NewChip(chip)
		defer c.Close()
		lines = c.Lines()
		reg.lines[chip] = lines
	} else {
		lines = l
	}
	return offset > 0 && offset < lines
}

func requestItem(chip string, offset int, opts ...ItemOption) (*item, error) {
	// check if chip exists
	if !isChip(chip) {
		return nil, fmt.Errorf("chip %s does not exist", chip)
	}
	// check if offset is valid
	if !isOffset(chip, offset) {
		return nil, fmt.Errorf("offset %d is out of range of chip %s, 0-%d", offset, chip, reg.lines[chip])
	}

	// apply options
	options := &ItemOptions{}
	for _, io := range opts {
		err := io.applyItemOption(options)
		if err != nil {
			return nil, err
		}
	}

	i, err := reg.Get(chip, offset)

	if _, ok := err.(ItemNotFound); ok {
		// item doesnt exist in registry so we'll create it
		i = &item{
			Line:     nil,
			RWMutex:  &sync.RWMutex{},
			chipName: chip,
			state:    options.state,
			offset:   offset,
			events:   newEventRegistry(),
			options:  options,
			tgc:      nil,
			// closed:   false,
		}

		t, err := tgc.New(i.tgcHandler)
		if err != nil {
			return nil, err
		}
		i.tgc = t

		err = reg.Add(chip, offset, i)
		if err != nil {
			return nil, err
		}

		logger.Infof("item registerd on line %d of chip %s as %s", offset, chip, options.mode)
	} else {
		// already exits, check if its of the same line direction
		i.RLock()
		info, err := i.Line.Info()
		i.RUnlock()
		if err != nil {
			return nil, err
		}
		if info.Config.Direction != gpiod.LineDirection(options.mode) {
			return nil, fmt.Errorf("this item is already registered as %v, you can't register it as %v", Mode(info.Config.Direction), options.mode)
		}
		logger.Infof("item registerd on line %d of chip %s as %s got a new owner", offset, chip, options.mode)
	}

	// either if it just got created or it was there all along, increment it's tgc
	i.RLock()
	tgc := i.tgc
	i.RUnlock()
	tgc.Add()

	return i, nil
}

// func GetChip(chip string) (c Chip)

func RequestItem(chip string, offset int, opts ...ItemOption) (Item, error) {
	return requestItem(chip, offset, opts...)
}

func RequestItemHelper(chip string, offset int, mode string, pull string, state string) (Item, error) {
	var (
		m Mode
		p Pull
		s State
	)
	err := m.Set(mode)
	if err != nil {
		return nil, err
	}

	switch m {
	case ModeInput:
		err = p.Set(pull)
		if err != nil {
			return nil, err
		}
		i, err := requestItem(chip, offset, AsInput(p))
		if err != nil {
			return nil, err
		}
		return i, nil

	case ModeOutput:
		err = s.Set(state)
		if err != nil {
			return nil, err
		}
		i, err := requestItem(chip, offset, AsOutput(s))
		if err != nil {
			return nil, err
		}
		return i, nil
	}

	panic("this should not happen")
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

func GetItem(chip string, offset int) (Item, error) {
	return reg.Get(chip, offset)
}

func SetState(chipName string, offset int, state State) error {
	i, err := GetItem(chipName, offset)
	if err != nil {
		return err
	}
	err = i.SetState(state)
	if err != nil {
		return err
	}
	return nil
}

func Close() (err error) {
	reg.ForEach(func(chip string, offset int, item *item) {
		item.shutdown()
	})
	return
}

func Chips() []string {
	return gpiod.Chips()
}

func Info() (map[string]map[int]ItemInfo, error) {
	info := map[string]map[int]ItemInfo{}
	var gErr error

	reg.ForEach(func(chip string, offset int, item *item) {
		i, err := item.Info()
		if err != nil {
			multierr.Append(gErr, err)
		}
		info[chip][offset] = i
	})
	if gErr != nil {
		return nil, gErr
	}

	return info, nil
}
