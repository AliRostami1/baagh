package tgc

import (
	"math"
	"sync"

	"github.com/AliRostami1/baagh/pkg/errprim"
)

type OnChangeCallback = func(bool)

type GarbageCollector interface {
	Count() int
	Add()
	Delete()
	State() bool
	Limit() int
}

type Tgc struct {
	ownerCount uint
	limit      uint
	state      bool
	onChange   OnChangeCallback
	*sync.RWMutex
}

func New(onChange OnChangeCallback, opts ...Option) *Tgc {
	options := &Options{
		limit:             math.MaxUint16,
		initialOwnerCount: 1,
	}
	for _, opt := range opts {
		opt.applyOption(options)
	}
	return &Tgc{
		ownerCount: options.initialOwnerCount,
		limit:      options.limit,
		state:      options.initialOwnerCount > 0,
		onChange:   onChange,
		RWMutex:    &sync.RWMutex{},
	}
}

func (t *Tgc) Count() uint {
	t.Lock()
	defer t.Unlock()
	return t.ownerCount
}

func (t *Tgc) Add() {
	t.Lock()
	if t.ownerCount == t.limit {
		return
	}
	t.ownerCount += 1
	t.Unlock()
	t.stateChangeCheck()
}

func (t *Tgc) Delete() {
	t.Lock()
	if t.ownerCount == 0 {
		return
	}
	t.ownerCount -= 1
	t.Unlock()
	t.stateChangeCheck()
}

func (t *Tgc) State() bool {
	t.Lock()
	defer t.Unlock()
	return t.state
}

func (t *Tgc) Limit() uint {
	t.Lock()
	defer t.Unlock()
	return t.limit
}

func (t *Tgc) stateChangeCheck() {
	t.Lock()
	if t.ownerCount > 0 != t.state {
		t.state = t.ownerCount > 0
		if t.onChange != nil {
			s := t.state
			t.Unlock()
			t.onChange(s)
		} else {
			t.Unlock()
		}
	} else {
		t.Unlock()
	}
}

type Option interface {
	applyOption(*Options) error
}

type Options struct {
	limit             uint
	initialOwnerCount uint
}

type LimitOption uint

func (l LimitOption) applyOption(o *Options) error {
	if l == 0 || uint(l) > uint(o.initialOwnerCount) {
		return errprim.OptionError{Field: "Limit", Value: l}
	}
	o.limit = uint(l)
	return nil
}

func WithLimit(limit uint) LimitOption {
	return LimitOption(limit)
}

type InitialOwnerCountOption uint

func (i InitialOwnerCountOption) applyOption(o *Options) error {
	if uint(i) > uint(o.limit) {
		return errprim.OptionError{Field: "Limit", Value: i}
	}
	o.initialOwnerCount = uint(i)
	return nil
}
