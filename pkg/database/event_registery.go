package database

import (
	"fmt"
	"sync"
)

type EventHandler func(key, value string)

type EventRegistery struct {
	registry map[string][]EventHandler
	*sync.RWMutex
}

func DefaultEventRegistry() *EventRegistery {
	return &EventRegistery{
		registry: make(map[string][]EventHandler),
		RWMutex:  &sync.RWMutex{},
	}
}

func (e *EventRegistery) addEvent(key string, fn ...EventHandler) {
	e.Lock()
	defer e.Unlock()
	e.registry[key] = append(e.registry[key], fn...)
}

func (e *EventRegistery) forEach(key string, fn func(fn EventHandler)) error {
	e.Lock()
	defer e.Unlock()
	eventHandlers, exists := e.registry[key]
	if !exists {
		return KeyNotFoundError{key}
	}
	for _, eh := range eventHandlers {
		fn(eh)
	}
	return nil
}

type KeyNotFoundError struct {
	key string
}

func (n KeyNotFoundError) Error() string {
	return fmt.Sprintf("key: %s not found", n.key)
}
