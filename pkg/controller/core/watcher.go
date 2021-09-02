package core

type Watcher interface {
	Closer
	Watch() <-chan *ItemEvent
	State() State
}

type watcher struct {
	item         *item
	eventChannel chan *ItemEvent
}

func (w *watcher) Watch() <-chan *ItemEvent {
	return w.eventChannel
}

func (w *watcher) Close() error {
	defer close(w.eventChannel)
	w.item.removeWatcher(w.eventChannel)
	return w.item.Close()
}

func (w *watcher) State() State {
	return w.item.State()
}
