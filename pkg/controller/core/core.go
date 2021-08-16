package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/logy"

	"github.com/warthog618/gpiod"
	"go.uber.org/multierr"
)

// key is chip name
var chips chipRegistry = chipRegistry{
	registry: map[string]*Chip{},
	RWMutex:  &sync.RWMutex{},
}

var events = eventRegistry{
	events:  []EventHandler{},
	RWMutex: &sync.RWMutex{},
}

var logger logy.Logger = logy.DummyLogger{}

func SetLogger(l logy.Logger) error {
	if l == nil {
		return fmt.Errorf("logger can't be nil")
	}
	logger = l
	return nil
}

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
		items: &itemRegistry{registry: map[int]*Item{}, RWMutex: &sync.RWMutex{}},
		mu:    &sync.RWMutex{},
	}
	err = chips.Append(options.name, chip)
	if err != nil {
		return nil, err
	}
	logger.Infof("chip %s registerd successfully by %s", options.name, options.consumer)
	return
}

func RegisterItem(chip string, offset int, opts ...ItemOption) (item *Item, err error) {
	// get the chip
	c, err := chips.Get(chip)
	if err != nil {
		return nil, fmt.Errorf("there is no registered chip named %s", chip)
	}

	return c.RegisterItem(offset, opts...)
}

func Subscribe(fns ...EventHandler) {
	events.AddEventListener(fns...)
}

func SetState(chipName string, offset int, state State) (err error) {
	c, err := chips.Get(chipName)
	if err != nil {
		return
	}
	i, err := c.GetItem(offset)
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
	c, err := chips.Get(chipName)
	if err != nil {
		return
	}
	i, err := c.items.Get(offset)
	if err != nil {
		return
	}
	return i.AddEventListener(fns...)
}

func Cleanup() (err error) {
	chips.ForEach(func(chipName string, chip *Chip) {
		err = multierr.Append(err, chip.Cleanup())
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		logger.Infof("gpio core is successfully cleanedup")
	}
	return
}

type Chip struct {
	chip  *gpiod.Chip
	items *itemRegistry

	mu *sync.RWMutex
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

	if _, ok := err.(ItemNotFound); !ok {
		// already exits, check if its of the same line direction
		info, err := item.line.Info()
		if err != nil {
			return nil, err
		}
		if info.Config.Direction != gpiod.LineDirection(options.io.mode) {
			return nil, fmt.Errorf("this item is already registered as %s", Mode(info.Config.Direction))
		}
		item.incrOwner()
		return item, nil
	}

	item = &Item{
		line:  nil,
		state: options.state,
		events: &eventRegistry{
			events:  []EventHandler{},
			RWMutex: &sync.RWMutex{},
		},
		ownerCount: 0,
		mu:         &sync.RWMutex{},
	}
	item.AddEventListener(func(event *ItemEvent) {
		events.CallAll(event)
	})

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
		return nil, fmt.Errorf("you have to set the mode")
	}

	err = c.items.Add(offset, item)
	if err != nil {
		return nil, err
	}
	logger.Infof("item registerd on line %o as %s", offset, options.io.mode)
	return
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

type Item struct {
	line       *gpiod.Line
	state      State
	events     *eventRegistry
	ownerCount int

	mu *sync.RWMutex
}

func (i *Item) Unregister() {
	i.decrOwner()
}

func (i *Item) incrOwner() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ownerCount += 1
	if i.ownerCount == 0 {
		i.Cleanup()
	}
}

func (i *Item) decrOwner() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ownerCount -= 1
	if i.ownerCount == 0 {
		i.Cleanup()
	}
}

func (i *Item) SetState(state State) (err error) {
	i.mu.Lock()
	iState := i.state
	line := i.line
	i.mu.Unlock()
	if iState == state {
		return
	}
	info, err := line.Info()
	if err != nil {
		return
	}
	if info.Config.Direction == gpiod.LineDirectionOutput {
		err = line.SetValue(int(state))
		if err != nil {
			return
		}
	}
	i.mu.Lock()
	i.state = state
	itemEvents := i.events
	i.mu.Unlock()

	events.CallAll(&ItemEvent{
		Item: i,
	})
	itemEvents.CallAll(&ItemEvent{
		Item: i,
	})
	logger.Debugf("state changed to %s on line %o of chip %s", state, i.line.Offset(), i.line.Chip())
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
	err = i.events.AddEventListener(fns...)
	return
}

func (i *Item) Cleanup() (err error) {
	i.mu.Lock()
	line := i.line
	i.mu.Unlock()
	c, err := GetChip(line.Chip())
	if err != nil {
		return
	}
	c.items.Delete(line.Offset())
	line.SetValue(int(Inactive))
	line.Close()
	i = nil
	logger.Infof("cleaned up item %o of %s", line.Offset(), line.Chip())
	return
}
