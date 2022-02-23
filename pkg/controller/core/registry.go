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
	i.Lock()
	reg := i.chips
	i.Unlock()

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
	i.Lock()
	reg := i.chips
	i.Unlock()
	delete(reg[chip], offset)
}

func (i *registry) Get(chip string, offset int) (*item, error) {
	i.Lock()
	reg := i.chips
	i.Unlock()

	item, ok := reg[chip][offset]
	if !ok {
		return nil, ItemNotFound{offset: offset}
	}
	return item, nil
}

func (i *registry) ForEach(chip string, fn func(offset int, item *item)) {
	i.Lock()
	reg := i.chips
	i.Unlock()
	for index, item := range reg[chip] {
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
