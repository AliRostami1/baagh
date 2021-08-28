package core

import (
	"fmt"
	"sync"

	"github.com/warthog618/gpiod"
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

type Item struct {
	line       *gpiod.Line
	state      State
	events     *eventRegistry
	ownerCount int

	mu *sync.RWMutex
}

func newItem(state State) *Item {
	return &Item{
		line:  nil,
		state: state,
		events: &eventRegistry{
			events:  []EventHandler{},
			RWMutex: &sync.RWMutex{},
		},
		ownerCount: 0,
		mu:         &sync.RWMutex{},
	}
}

func (i *Item) Unregister() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ownerCount -= 1
	if i.ownerCount == 0 {
		i.Cleanup()
	}
}

func (i *Item) incrOwner() {
	i.mu.Lock()
	i.ownerCount += 1
	i.mu.Unlock()
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
