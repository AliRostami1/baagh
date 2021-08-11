package chip

import (
	"fmt"
	"sync"
)

type objectRegistry struct {
	registry map[int]*Object
	*sync.RWMutex
}

func defaultObjectRegistry() *objectRegistry {
	return &objectRegistry{
		registry: map[int]*Object{},
		RWMutex:  &sync.RWMutex{},
	}
}

func (o *objectRegistry) append(pin int, item *Object) error {
	o.Lock()
	defer o.Unlock()

	if _, ok := o.registry[pin]; ok {
		return AlreadyRegisteredError{Pin: pin}
	}
	o.registry[pin] = item
	return nil
}

func (o *objectRegistry) item(pin int) (*Object, error) {
	o.Lock()
	defer o.Unlock()

	object, ok := o.registry[pin]
	if !ok {
		return nil, KeyNotFoundError{Pin: pin}
	}
	return object, nil
}

func (o *objectRegistry) forEach(fn func(index int, item *Object)) {
	o.Lock()
	defer o.Unlock()
	for index, item := range o.registry {
		fn(index, item)
	}
}

type AlreadyRegisteredError struct {
	Pin int
}

func (a AlreadyRegisteredError) Error() string {
	return fmt.Sprintf("pin: %o is already registered", a.Pin)
}

type KeyNotFoundError struct {
	Pin int
}

func (n KeyNotFoundError) Error() string {
	return fmt.Sprintf("there is no controller with pin number: %o", n.Pin)
}
