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

func New(onChange OnChangeCallback, opts ...Option) (tgc *Tgc, err error) {
	if onChange == nil {
		return nil, errprim.OptionError{Field: "onChange", Value: onChange}
	}
	options := &Options{
		limit:             math.MaxUint16,
		initialOwnerCount: 0,
	}
	for _, opt := range opts {
		err = opt.applyOption(options)
		if err != nil {
			return
		}
	}
	tgc = &Tgc{
		ownerCount: options.initialOwnerCount,
		limit:      options.limit,
		state:      false,
		onChange:   onChange,
		RWMutex:    &sync.RWMutex{},
	}
	tgc.stateChangeCheck()
	return
}

func (t *Tgc) Add() {
	t.Lock()
	if t.ownerCount == t.limit {
		t.Unlock()
		return
	}
	t.ownerCount += 1
	t.Unlock()
	t.stateChangeCheck()
}

func (t *Tgc) Delete() {
	t.Lock()
	if t.ownerCount == 0 {
		t.Unlock()
		return
	}
	t.ownerCount -= 1
	t.Unlock()
	t.stateChangeCheck()
}

func (t *Tgc) Shutdown() {
	t.Lock()
	if t.ownerCount == 0 {
		t.Unlock()
		return
	}
	t.ownerCount = 0
	t.Unlock()
	t.stateChangeCheck()
}

func (t *Tgc) Count() uint {
	t.RLock()
	defer t.RUnlock()
	return t.ownerCount
}

func (t *Tgc) State() bool {
	t.RLock()
	defer t.RUnlock()
	return t.state
}

func (t *Tgc) Limit() uint {
	t.RLock()
	defer t.RUnlock()
	return t.limit
}

func (t *Tgc) stateChangeCheck() {
	t.Lock()
	if (t.ownerCount > 0) != t.state {
		t.state = t.ownerCount > 0
		s := t.state
		t.Unlock()
		t.onChange(s)
		return
	}
	t.Unlock()
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
	if l == 0 || uint(l) < uint(o.initialOwnerCount) {
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

func WithInitialOwnerCount(cnt uint) InitialOwnerCountOption {
	return InitialOwnerCountOption(cnt)
}
