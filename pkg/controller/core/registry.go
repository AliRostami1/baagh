package core

import (
	"fmt"
	"sync"
)

type chipRegistry struct {
	registry map[string]*Chip
	*sync.RWMutex
}

func (c *chipRegistry) append(name string, chip *Chip) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.registry[name]; ok {
		return DuplicateChipError{Chip: name}
	}
	c.registry[name] = chip
	return nil
}

func (c *chipRegistry) get(name string) (*Chip, error) {
	c.Lock()
	defer c.Unlock()

	chip, ok := c.registry[name]
	if !ok {
		return nil, ChipNotFoundError{Chip: name}
	}
	return chip, nil
}

func (c *chipRegistry) forEach(fn func(chipName string, item *Chip)) {
	c.Lock()
	defer c.Unlock()
	for index, item := range c.registry {
		fn(index, item)
	}
}

type DuplicateChipError struct {
	Chip string
}

func (d DuplicateChipError) Error() string {
	return fmt.Sprintf("chip: %s is already registered", d.Chip)
}

type ChipNotFoundError struct {
	Chip string
}

func (c ChipNotFoundError) Error() string {
	return fmt.Sprintf("there is no chip with named: %s", c.Chip)
}

type itemRegistry struct {
	registry map[int]*Item
	*sync.RWMutex
}

func (i *itemRegistry) append(pin int, item *Item) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.registry[pin]; ok {
		return DuplicateItemError{Pin: pin}
	}
	i.registry[pin] = item
	return nil
}

func (i *itemRegistry) get(pin int) (*Item, error) {
	i.Lock()
	defer i.Unlock()

	item, ok := i.registry[pin]
	if !ok {
		return nil, ItemNotFound{Pin: pin}
	}
	return item, nil
}

func (i *itemRegistry) forEach(fn func(pin int, item *Item)) {
	i.Lock()
	defer i.Unlock()
	for index, item := range i.registry {
		fn(index, item)
	}
}

type DuplicateItemError struct {
	Pin int
}

func (a DuplicateItemError) Error() string {
	return fmt.Sprintf("pin: %o is already registered", a.Pin)
}

type ItemNotFound struct {
	Pin int
}

func (n ItemNotFound) Error() string {
	return fmt.Sprintf("there is no item registered on pin: %o", n.Pin)
}

type ItemEvent struct {
}

type EventHandler func(event *ItemEvent)

type EventRegistry struct {
	events []EventHandler
	*sync.RWMutex
}

func (o *EventRegistry) addEventListener(fn ...EventHandler) error {
	o.Lock()
	defer o.Unlock()

	o.events = append(o.events, fn...)
	return nil
}

func (o *EventRegistry) forEach(cb func(index int, handler EventHandler)) {
	o.Lock()
	defer o.Unlock()
	for index, eh := range o.events {
		cb(index, eh)
	}
}
