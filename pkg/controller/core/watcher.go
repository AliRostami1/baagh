package core

import (
	"go.uber.org/multierr"
)

type ItemEvent struct {
	ItemInfo
}

type Watcher interface {
	Closer
	Watch() <-chan *ItemEvent
}

type watcher struct {
	item         *item
	chip         *chip
	eventChannel chan *ItemEvent
}

func (w *watcher) Watch() <-chan *ItemEvent {
	return w.eventChannel
}

func (w *watcher) Close() error {
	defer close(w.eventChannel)
	w.item.removeWatcher(w.eventChannel)

	return multierr.Combine(w.item.Close(), w.chip.Close())
}
