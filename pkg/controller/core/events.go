package core

import "sync"

type EventChannel = <-chan *ItemEvent

type eventRegistry struct {
	events []chan *ItemEvent
	*sync.RWMutex
}

func newEventRegistry() *eventRegistry {
	return &eventRegistry{
		events:  []chan *ItemEvent{},
		RWMutex: &sync.RWMutex{},
	}
}

func (e *eventRegistry) Add(ch chan *ItemEvent) {
	e.Lock()
	defer e.Unlock()
	e.events = append(e.events, ch)
}

func (e *eventRegistry) Remove(ch chan *ItemEvent) {
	e.Lock()
	defer e.Unlock()
	for i, channel := range e.events {
		if channel == ch {
			e.events[i] = e.events[len(e.events)-1]
			e.events = e.events[:len(e.events)-1]
		}
	}
}

func (e *eventRegistry) ForEach(cb func(index int, ch chan *ItemEvent)) {
	e.Lock()
	defer e.Unlock()
	for index, ch := range e.events {
		cb(index, ch)
	}
}

func (e *eventRegistry) CallAll(evt *ItemEvent) {
	go func() {
		e.Lock()
		events := e.events
		e.Unlock()
		for _, ch := range events {
			ch <- evt
		}
	}()
}
