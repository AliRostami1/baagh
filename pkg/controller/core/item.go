package core

import (
	"sync"

	"github.com/warthog618/gpiod"
)

type ItemI interface {
	Register() error
	Unregister()
	Active() bool
	SetState() error
	State() State
	Mode() Mode
	Info()
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
		i.Close()
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

func (i *Item) Close() (err error) {
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
