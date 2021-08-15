package security

import (
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/core"
)

type registry struct {
	registry map[string]*Security
	*sync.RWMutex
}

func (r *registry) Append(tag string, security *Security) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.registry[tag]; ok {
		return DuplicateTagError{Tag: tag}
	}
	r.registry[tag] = security
	return nil
}

func (r *registry) Get(tag string) (*Security, error) {
	r.Lock()
	defer r.Unlock()

	chip, ok := r.registry[tag]
	if !ok {
		return nil, TagNotFoundError{Tag: tag}
	}
	return chip, nil
}

func (r *registry) ForEach(fn func(chipName string, security *Security)) {
	r.Lock()
	defer r.Unlock()
	for index, chip := range r.registry {
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
	registry map[string]map[int]*core.Item
	*sync.RWMutex
}

func (i *itemRegistry) Add(chip string, offset int, item *core.Item) error {
	i.Lock()
	defer i.Unlock()
	if _, ok := i.registry[chip][offset]; ok {
		return DuplicateItemError{Chip: chip, Offset: offset}
	}
	i.registry[chip][offset] = item
	return nil
}

func (i *itemRegistry) Get(chip string, offset int) (*core.Item, error) {
	i.Lock()
	defer i.Unlock()
	if item, ok := i.registry[chip][offset]; ok {
		return item, nil
	}
	return nil, ItemNotFoundError{Chip: chip, Offset: offset}
}

func (i *itemRegistry) ForEach(fn func(i *core.Item)) {
	i.Lock()
	defer i.Unlock()
	for _, c := range i.registry {
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
