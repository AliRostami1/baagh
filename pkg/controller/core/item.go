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
	Closer
	SetState(State) error
	State() State
	Info() (*ItemInfo, error)
}

type item struct {
	*gpiod.Line
	*sync.RWMutex
	chip    *chip
	state   State
	offset  int
	events  *eventRegistry
	options *ItemOptions
	tgc     *tgc.Tgc
}

func (i *item) tgcHandler(b bool) {
	// TODO: we are ignoring error here, fix it!
	i.Lock()
	offset, chip, options := i.offset, i.chip, i.options
	i.Unlock()
	if b {
		switch options.mode {
		case ModeInput:
			var (
				l   *gpiod.Line
				err error
			)
			switch options.pull {
			case PullDisabled, PullUnknown:
				l, err = chip.RequestLine(offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges)
			case PullDown:
				l, err = chip.RequestLine(offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges, gpiod.WithPullDown)
			case PullUp:
				l, err = chip.RequestLine(offset, gpiod.AsInput, gpiod.WithEventHandler(i.eventHandler), gpiod.WithBothEdges, gpiod.WithPullUp)
			}
			if err != nil {
				logger.Errorf("requestLine failed: %v", err)
			}
			i.Lock()
			i.Line = l
			i.Unlock()
		case ModeOutput:
			l, err := chip.RequestLine(offset, gpiod.AsOutput(int(options.state)))
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
		chip.Lock()
		chip.items.Delete(offset)
		chip.Unlock()
		i.cleanup()
		return
	}
}

func (i *item) eventHandler(evt gpiod.LineEvent) {
	var newState State
	switch evt.Type {
	case gpiod.LineEventRisingEdge:
		newState = StateActive
	case gpiod.LineEventFallingEdge:
		newState = StateInactive
	}
	i.Lock()
	if i.state == newState {
		i.Unlock()
		return
	}
	i.Unlock()
	i.setState(newState)
	i.eventEmmiter(&evt)
}

func (i *item) eventEmmiter(evt *gpiod.LineEvent) error {
	i.Lock()
	itemEvents := i.events
	i.Unlock()
	info, err := i.Info()
	if err != nil {
		return err
	}

	// events.CallAll(&ItemEvent{
	// 	item: i,
	// })
	itemEvents.CallAll(&ItemEvent{
		Info:        info,
		Item:        i,
		IsLineEvent: evt != nil,
		LineEvent:   evt,
	})
	return nil
}

func (i *item) SetState(state State) (err error) {
	i.Lock()
	if i.state == state {
		i.Unlock()
		return
	}
	i.Unlock()

	err = i.setState(state)
	if err != nil {
		return
	}
	return i.eventEmmiter(nil)
}

func (i *item) setState(state State) (err error) {
	i.Lock()
	line := i.Line
	i.Unlock()

	info, err := i.Info()
	if err != nil {
		return
	}
	if info.Config.Direction == gpiod.LineDirectionOutput {
		err = line.SetValue(int(state))
		if err != nil {
			return
		}
	}
	i.Lock()
	i.state = state
	i.Unlock()

	logger.Debugf("state changed to %s on line %o of chip %s", state, line.Offset(), line.Chip())
	return
}

func (i *item) State() State {
	i.Lock()
	defer i.Unlock()
	return i.state
}

func (i *item) Info() (*ItemInfo, error) {
	i.Lock()
	li, err := i.Line.Info()
	i.Unlock()
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
	i.Lock()
	defer i.Unlock()
	return i.options.mode
}

func (i *item) Pull() Pull {
	i.Lock()
	defer i.Unlock()
	return i.options.pull
}

func (i *item) removeWatcher(ch chan *ItemEvent) {
	i.Lock()
	ev := i.events
	i.Unlock()
	ev.Remove(ch)
}

func (i *item) Close() error {
	i.Lock()
	tgc := i.tgc
	chip := i.chip
	i.Unlock()
	tgc.Delete()
	chip.Close()
	return nil
}

func (i *item) cleanup() (err error) {
	i.Lock()
	line := i.Line
	c := i.chip
	i.Unlock()

	c.items.Delete(line.Offset())
	i.setState(StateInactive)
	line.Close()
	i = nil
	logger.Infof("cleaned up item %o of %s", line.Offset(), line.Chip())
	return
}
