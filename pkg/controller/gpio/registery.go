package gpio

import (
	"fmt"
	"sync"
)

type ItemRegistry struct {
	registry map[string]*Object
	*sync.RWMutex
}

func DefaultItemRegistry() *ItemRegistry {
	return &ItemRegistry{
		registry: make(map[string]*Object),
		RWMutex:  &sync.RWMutex{},
	}
}

func (i *ItemRegistry) addItem(pin int, item *Object) error {
	i.Lock()
	defer i.Unlock()

	_, exists := i.registry[makeKey(pin)]
	if exists {
		return &AlreadyRegisteredError{pin: pin}
	}
	i.registry[makeKey(pin)] = item
	return nil
}

func (i *ItemRegistry) getItem(pin int) (*Object, error) {
	i.Lock()
	defer i.Unlock()

	item, exists := i.registry[makeKey(pin)]
	if !exists {
		return nil, KeyNotFoundError{pin: pin, key: makeKey(pin)}
	}
	return item, nil
}

func (i *ItemRegistry) forEach(fn func(item *Object)) {
	i.Lock()
	defer i.Unlock()
	for _, item := range i.registry {
		fn(item)
	}
}

type AlreadyRegisteredError struct {
	pin int
}

func (a *AlreadyRegisteredError) Error() string {
	return fmt.Sprintf("pin: %o is already registered", a.pin)
}

type KeyNotFoundError struct {
	pin int
	key string
}

func (n KeyNotFoundError) Error() string {
	return fmt.Sprintf("there is no controller with pin number: %o, and corresponding key: %s", n.pin, n.key)
}
