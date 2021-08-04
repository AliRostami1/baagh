package gpio

import (
	"fmt"
	"sync"
)

type ItemRegistery struct {
	registry map[string]*Item
	*sync.RWMutex
}

func DefaultItemRegistry() *ItemRegistery {
	return &ItemRegistery{
		registry: make(map[string]*Item),
		RWMutex:  &sync.RWMutex{},
	}
}

func (i *ItemRegistery) addItem(pin uint8, item *Item) error {
	i.Lock()
	defer i.Unlock()

	_, exists := i.registry[makeKey(pin)]
	if exists {
		return &AlreadyRegisteredError{pin: pin}
	}
	i.registry[makeKey(pin)] = item
	return nil
}

func (i *ItemRegistery) getItem(pin uint8) (*Item, error) {
	i.Lock()
	defer i.Unlock()

	item, exists := i.registry[makeKey(pin)]
	if !exists {
		return nil, KeyNotFoundError{pin: pin, key: makeKey(pin)}
	}
	return item, nil
}

func (i *ItemRegistery) forEach(fn func(item *Item)) {
	i.Lock()
	defer i.Unlock()
	for _, item := range i.registry {
		fn(item)
	}
}

type AlreadyRegisteredError struct {
	pin uint8
}

func (a *AlreadyRegisteredError) Error() string {
	return fmt.Sprintf("pin: %o is already registered", a.pin)
}

type KeyNotFoundError struct {
	pin uint8
	key string
}

func (n KeyNotFoundError) Error() string {
	return fmt.Sprintf("there is no controller with pin number: %o, and corresponding key: %s", n.pin, n.key)
}
