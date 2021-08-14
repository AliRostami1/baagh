package core

import (
	"fmt"
	"sync"
)

type chipRegistry struct {
	registry map[string]*Chip
	*sync.RWMutex
}

func (c *chipRegistry) Append(name string, chip *Chip) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.registry[name]; ok {
		return DuplicateChipError{Chip: name}
	}
	c.registry[name] = chip
	return nil
}

func (c *chipRegistry) Get(name string) (*Chip, error) {
	c.Lock()
	defer c.Unlock()

	chip, ok := c.registry[name]
	if !ok {
		return nil, ChipNotFoundError{Chip: name}
	}
	return chip, nil
}

func (c *chipRegistry) ForEach(fn func(chipName string, chip *Chip)) {
	c.Lock()
	defer c.Unlock()
	for index, chip := range c.registry {
		fn(index, chip)
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

func (i *itemRegistry) Append(offset int, item *Item) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.registry[offset]; ok {
		return DuplicateItemError{offset: offset}
	}
	i.registry[offset] = item
	return nil
}

func (i *itemRegistry) Get(offset int) (*Item, error) {
	i.Lock()
	defer i.Unlock()

	item, ok := i.registry[offset]
	if !ok {
		return nil, ItemNotFound{offset: offset}
	}
	return item, nil
}

func (i *itemRegistry) ForEach(fn func(offset int, item *Item)) {
	i.Lock()
	defer i.Unlock()
	for index, item := range i.registry {
		fn(index, item)
	}
}

type DuplicateItemError struct {
	offset int
}

func (a DuplicateItemError) Error() string {
	return fmt.Sprintf("offset: %o is already registered", a.offset)
}

type ItemNotFound struct {
	offset int
}

func (n ItemNotFound) Error() string {
	return fmt.Sprintf("there is no item registered on offset: %o", n.offset)
}

type ItemEvent struct {
	Item *Item
}

type EventHandler func(event *ItemEvent)

type eventRegistry struct {
	events []EventHandler
	*sync.RWMutex
}

func (e *eventRegistry) AddEventListener(fn ...EventHandler) error {
	e.Lock()
	defer e.Unlock()

	e.events = append(e.events, fn...)
	return nil
}

func (e *eventRegistry) ForEach(cb func(index int, handler EventHandler)) {
	e.Lock()
	defer e.Unlock()
	for index, eh := range e.events {
		cb(index, eh)
	}
}

func (e *eventRegistry) CallAll(evt *ItemEvent) {
	go func() {
		e.Lock()
		defer e.Unlock()
		for _, eh := range e.events {
			eh(evt)
		}
	}()
}
