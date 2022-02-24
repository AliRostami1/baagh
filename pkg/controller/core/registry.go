package core

import (
	"fmt"
	"sync"
)

type registry struct {
	chips map[string]map[int]*item
	lines map[string]int
	*sync.RWMutex
}

func (i *registry) Add(chip string, offset int, it *item) error {
	i.RLock()
	reg := i.chips
	i.RUnlock()

	if _, ok := reg[chip]; !ok {
		reg[chip] = map[int]*item{}
	}
	if _, ok := reg[chip][offset]; ok {
		return DuplicateItemError{offset: offset}
	}

	reg[chip][offset] = it
	return nil
}

func (i *registry) Delete(chip string, offset int) {
	i.RLock()
	reg := i.chips
	i.RUnlock()
	delete(reg[chip], offset)
}

func (i *registry) Get(chip string, offset int) (*item, error) {
	i.RLock()
	reg := i.chips
	i.RUnlock()

	if _, ok := reg[chip]; !ok {
		reg[chip] = map[int]*item{}
		return nil, ItemNotFound{offset: offset}
	}

	item, ok := reg[chip][offset]
	if !ok {
		return nil, ItemNotFound{offset: offset}
	}
	return item, nil
}

func (i *registry) ForEach(chip string, fn func(offset int, item *item)) {
	i.RLock()
	reg := i.chips
	i.RUnlock()
	c, ok := reg[chip]
	if !ok {
		return
	}
	for index, item := range c {
		fn(index, item)
	}
}

func (i *registry) GetAll(chip string) (map[int]*item, error) {
	i.RLock()
	reg := i.chips
	i.RUnlock()
	c, ok := reg[chip]
	if !ok {
		return nil, fmt.Errorf("chip not found")
	}
	return c, nil
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
