package core

import (
	"sync"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/warthog618/gpiod"
)

type ItemInfo struct {
}

type Item interface {
	Closer
	NewWatcher() (Watcher, error)
	SetState(State) error
	State() State
	Info() ItemInfo
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
	offset := i.offset
	chip := i.chip
	defer i.Unlock()
	options := i.options
	if b {
		switch options.mode {
		case Input:
			handler := func(evt gpiod.LineEvent) {
				switch evt.Type {
				case gpiod.LineEventRisingEdge:
					i.SetState(Active)
				case gpiod.LineEventFallingEdge:
					i.SetState(Inactive)
				}
			}
			l, _ := chip.RequestLine(offset, gpiod.AsInput, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
			// if err != nil {
			// 	return
			// }
			i.Line = l
		case Output:
			l, _ := chip.RequestLine(offset, gpiod.AsOutput(int(options.state)))
			// if err != nil {
			// 	return nil, err
			// }
			i.Line = l
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

func (i *item) SetState(state State) (err error) {
	i.Lock()
	iState := i.state
	line := i.Line
	i.Unlock()
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
	i.Lock()
	i.state = state
	itemEvents := i.events
	i.Unlock()

	// events.CallAll(&ItemEvent{
	// 	item: i,
	// })
	itemEvents.CallAll(ItemEvent{
		// LineEvent: i,
	})
	logger.Debugf("state changed to %s on line %o of chip %s", state, line.Offset(), line.Chip())
	return
}

func (i *item) State() State {
	i.Lock()
	defer i.Unlock()
	return i.state
}

func (i *item) Info() ItemInfo {
	return ItemInfo{}
}

func (i *item) NewWatcher() (Watcher, error) {
	chip, err := getChip(i.Chip())
	if err != nil {
		return nil, err
	}

	w := &watcher{
		item:         i,
		chip:         chip,
		eventChannel: make(chan ItemEvent),
	}

	i.Lock()
	ev := i.events
	i.Unlock()
	ev.Add(w.eventChannel)

	return w, nil
}

func (i *item) removeWatcher(ch EventChannel) {
	i.Lock()
	ev := i.events
	i.Unlock()
	ev.Remove(ch)
}

func (i *item) Close() error {
	i.Lock()
	i.tgc.Delete()
	i.Unlock()
	return nil
}

func (i *item) cleanup() (err error) {
	i.Lock()
	line := i.Line
	i.Unlock()
	c, err := getChip(line.Chip())
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
