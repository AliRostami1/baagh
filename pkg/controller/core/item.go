package core

import (
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
)

type ItemInfo struct {
	*gpiod.LineInfo
	State State
	Mode  Mode
	Pull  Pull
}

type Item interface {
	Close() error
	SetState(State) error
	State() State
	Offset() int
	Mode() Mode
	Pull() Pull
	Chip() string
	Info() (*ItemInfo, error)
}

type item struct {
	*gpiod.Line
	*sync.RWMutex
	chipName string
	state    State
	offset   int
	events   *eventRegistry
	options  *ItemOptions
	tgc      *tgc.Tgc
	closed   bool
}

func (i *item) tgcHandler(b bool) {
	// TODO: we are ignoring error here, fix it!
	i.RLock()
	offset, options, chipName := i.offset, i.options, i.chipName
	i.RUnlock()
	if b {
		switch options.mode {
		case ModeInput:
			var (
				l   *gpiod.Line
				err error
			)
			switch options.pull {
			case PullDisabled, PullUnknown:
				l, err = gpiod.RequestLine(chipName, offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges)
			case PullDown:
				l, err = gpiod.RequestLine(chipName, offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges, gpiod.WithPullDown)
			case PullUp:
				l, err = gpiod.RequestLine(chipName, offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges, gpiod.WithPullUp)
			}
			if err != nil {
				logger.Errorf("requestLine failed: %v", err)
			}
			i.Lock()
			i.Line = l
			i.Unlock()
		case ModeOutput:
			l, err := gpiod.RequestLine(chipName, offset, gpiod.AsOutput(int(options.state)))
			if err != nil {
				logger.Errorf("requestLine failed: %v", err)
			}
			i.Lock()
			i.Line = l
			i.Unlock()
		default:
			return
			// return nil, fmt.Errorf("you have to set the mode")
		}
	} else {
		i.cleanup()
		return
	}
}

func (i *item) eventHandler(evt gpiod.LineEvent) {
	if !i.closed {
		var newState State
		switch evt.Type {
		case gpiod.LineEventRisingEdge:
			newState = StateActive
		case gpiod.LineEventFallingEdge:
			newState = StateInactive
		}
		i.RLock()

		i.eventEmmiter(&evt)

		if i.state == newState {
			i.RUnlock()
			return
		}
		i.RUnlock()
		i.setState(newState)
	}
}

func (i *item) eventEmmiter(evt *gpiod.LineEvent) error {
	info, err := i.Info()
	if err != nil {
		return err
	}
	i.RLock()
	itemEvents := i.events
	i.RUnlock()

	itemEvents.CallAll(&ItemEvent{
		Info:        info,
		Item:        i,
		IsLineEvent: evt != nil,
		LineEvent:   evt,
	})
	return nil
}

func (i *item) setState(state State) (err error) {
	i.RLock()
	line := i.Line
	i.RUnlock()

	err = line.SetValue(int(state))
	if err != nil {
		return
	}

	i.Lock()
	i.state = state
	i.Unlock()

	logger.Debugf("state changed to %s on line %d of chip %s", state, line.Offset(), line.Chip())
	return
}

func (i *item) removeWatcher(ch chan *ItemEvent) {
	i.RLock()
	ev := i.events
	i.RUnlock()
	ev.Remove(ch)
}

func (i *item) cleanup() (err error) {
	// set it's state to inactive
	i.setState(StateInactive)

	i.RLock()
	i.closed = true
	line, events := i.Line, i.events
	i.RUnlock()

	// delete the item from chip's item registry
	// c.items.Delete(i.Offset())

	// clear all event channels
	events.Cleanup()

	// close the gpiod.Line
	err = line.Close()
	if err != nil {
		return
	}

	logger.Infof("item %d is successfully cleaned up", line.Offset())

	return nil
}

// Exported
func (i *item) SetState(state State) (err error) {
	i.RLock()
	if i.state == state {
		i.RUnlock()
		return
	}
	i.RUnlock()

	err = i.setState(state)
	if err != nil {
		return
	}
	return i.eventEmmiter(nil)
}

func (i *item) State() State {
	i.RLock()
	defer i.RUnlock()
	return i.state
}

func (i *item) Info() (*ItemInfo, error) {
	i.RLock()
	li, err := i.Line.Info()
	i.RUnlock()
	if err != nil {
		return nil, err
	}
	return &ItemInfo{
		LineInfo: &li,
		State:    i.State(),
		Mode:     i.Mode(),
		Pull:     i.Pull(),
	}, nil
}

func (i *item) Mode() Mode {
	i.RLock()
	defer i.RUnlock()
	return i.options.mode
}

func (i *item) Pull() Pull {
	i.RLock()
	defer i.RUnlock()
	return i.options.pull
}

func (i *item) Offset() int {
	i.RLock()
	defer i.RUnlock()
	return i.offset
}

func (i *item) Chip() string {
	i.RLock()
	defer i.RUnlock()
	return i.chipName
}

func (i *item) Close() error {
	i.RLock()
	tgc := i.tgc
	i.RUnlock()
	tgc.Delete()
	return nil
}
