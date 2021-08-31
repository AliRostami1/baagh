package general

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

type registry struct {
	registry map[string]*General
	*sync.RWMutex
}

func (r *registry) Append(tag string, General *General) error {
	r.Lock()
	reg := r.registry
	r.Unlock()

	if _, ok := reg[tag]; ok {
		return DuplicateTagError{Tag: tag}
	}
	reg[tag] = General
	return nil
}

func (r *registry) Get(tag string) (*General, error) {
	r.Lock()
	reg := r.registry
	r.Unlock()

	chip, ok := reg[tag]
	if !ok {
		return nil, TagNotFoundError{Tag: tag}
	}
	return chip, nil
}

func (r *registry) ForEach(fn func(chipName string, General *General)) {
	r.Lock()
	reg := r.registry
	r.Unlock()
	for index, chip := range reg {
		fn(index, chip)
	}
}

type DuplicateTagError struct {
	Tag string
}

func (d DuplicateTagError) Error() string {
	return fmt.Sprintf("tag \"%s\" is already registered", d.Tag)
}

type TagNotFoundError struct {
	Tag string
}

func (t TagNotFoundError) Error() string {
	return fmt.Sprintf("there is no tag named \"%s\"", t.Tag)
}

type itemRegistry struct {
	registry map[string]map[int]core.Item
	*sync.RWMutex
}

func (i *itemRegistry) Add(chip string, offset int, item core.Item) error {
	i.Lock()
	reg := i.registry
	i.Unlock()
	if _, ok := reg[chip]; !ok {
		reg[chip] = make(map[int]core.Item)
	}
	reg[chip][offset] = item
	return nil
}

func (i *itemRegistry) Get(chip string, offset int) (core.Item, error) {
	i.Lock()
	reg := i.registry
	i.Unlock()
	if _, ok := reg[chip]; !ok {
		return nil, ItemNotFoundError{Chip: chip, Offset: offset}
	}
	if item, ok := reg[chip][offset]; ok {
		return item, nil
	}
	return nil, ItemNotFoundError{Chip: chip, Offset: offset}
}

func (i *itemRegistry) ForEach(fn func(i core.Item)) {
	i.Lock()
	reg := i.registry
	i.Unlock()
	for _, c := range reg {
		for _, i := range c {
			fn(i)
		}
	}
}

func (i *itemRegistry) Delete(chip string, offset int) {
	i.Lock()
	defer i.Unlock()
	delete(i.registry[chip], offset)
}

type ItemNotFoundError struct {
	Chip   string
	Offset int
}

func (i ItemNotFoundError) Error() string {
	return fmt.Sprintf("there is no item with %o offset on chip %s", i.Offset, i.Chip)
}

type DuplicateItemError struct {
	Chip   string
	Offset int
}

func (d DuplicateItemError) Error() string {
	return fmt.Sprintf("item with %o offset is already registered on chip %s", d.Offset, d.Chip)
}
