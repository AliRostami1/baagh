package core

type Watcher interface {
	Close() error
	Watch() <-chan *ItemEvent
	State() State
	Closed() bool
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

func (w *watcher) Closed() bool {
	return w.item.Closed()
}
