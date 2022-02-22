package core

import (
	"sync"

	"github.com/warthog618/gpiod"
)

type ItemEvent struct {
	Info        *ItemInfo
	Item        Item
	IsLineEvent bool
	*gpiod.LineEvent
}

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

func (e *eventRegistry) Cleanup() {
	e.Lock()
	defer e.Unlock()
	for _, ch := range e.events {
		close(ch)
	}
	e.events = e.events[:0]
}

func (e *eventRegistry) CallAll(evt *ItemEvent) {
	go func() {
		e.Lock()
		// events := e.events
		for _, ch := range e.events {
			ch <- evt
		}
		e.Unlock()
	}()
}
