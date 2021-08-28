package tgc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNew(t *testing.T) {
	called := 0
	cf := func(b bool) {
		called += 1
	}
	tgc, _ := New(cf)

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(0), tgc.ownerCount)
		assert.Equal(t, uint(math.MaxUint16), tgc.limit)
		assert.Equal(t, false, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 0, called)
	}

	for i := 0; i < 10; i += 1 {
		tgc.Add()
	}

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(10), tgc.ownerCount)
		assert.Equal(t, uint(math.MaxUint16), tgc.limit)
		assert.Equal(t, true, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 1, called)
	}

	for i := 0; i < 20; i += 1 {
		tgc.Delete()
	}

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(0), tgc.ownerCount)
		assert.Equal(t, uint(math.MaxUint16), tgc.limit)
		assert.Equal(t, false, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 2, called)
	}
}

func TestNewCustom(t *testing.T) {
	const (
		limit     = 10
		initCount = 5
	)

	called := 0
	cf := func(b bool) {
		called += 1
	}
	tgc, _ := New(cf, WithLimit(limit), WithInitialOwnerCount(initCount))

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(initCount), tgc.ownerCount)
		assert.Equal(t, uint(limit), tgc.limit)
		assert.Equal(t, true, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 1, called)
	}

	for i := 0; i < 10; i += 1 {
		tgc.Add()
	}

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(limit), tgc.ownerCount)
		assert.Equal(t, uint(limit), tgc.limit)
		assert.Equal(t, true, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 1, called)
	}

	for i := 0; i < 20; i += 1 {
		tgc.Delete()
	}

	if assert.NotNil(t, tgc) {
		assert.Equal(t, uint(0), tgc.ownerCount)
		assert.Equal(t, uint(limit), tgc.limit)
		assert.Equal(t, false, tgc.state)
		assert.NotNil(t, tgc.onChange)

		assert.Equal(t, tgc.Count(), tgc.ownerCount)
		assert.Equal(t, tgc.Limit(), tgc.limit)
		assert.Equal(t, tgc.State(), tgc.state)

		assert.Equal(t, 2, called)
	}
}

func TestStateCheck(t *testing.T) {
	trg, _ := New(func(b bool) {})

	if assert.NotNil(t, trg) {
		assert.False(t, trg.state)
	}

	trg.ownerCount += 1
	trg.stateChangeCheck()

	assert.True(t, trg.state)
}
