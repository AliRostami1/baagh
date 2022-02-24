package tgc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultNew(t *testing.T) {
	called := 0
	cf := func(b bool) {
		called += 1
	}
	tgc, err := New(cf)
	assert.Nil(t, err)
	require.NotNil(t, tgc)

	// empty tgc
	assert.Equal(t, uint(0), tgc.ownerCount)
	assert.Equal(t, uint(math.MaxUint16), tgc.limit)
	assert.False(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 0, called)

	// add 1
	tgc.Add()

	assert.Equal(t, uint(1), tgc.ownerCount)
	assert.Equal(t, uint(math.MaxUint16), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// add 9, with prev 1 = 10
	for i := 0; i < 9; i += 1 {
		tgc.Add()
	}

	assert.Equal(t, uint(10), tgc.ownerCount)
	assert.Equal(t, uint(math.MaxUint16), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// remove 1, with prev 10 = 9
	tgc.Delete()

	assert.Equal(t, uint(9), tgc.ownerCount)
	assert.Equal(t, uint(math.MaxUint16), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// remove 20, with prev 9 = 0 (it shouldn't go below 0)
	for i := 0; i < 20; i += 1 {
		tgc.Delete()
	}

	assert.Equal(t, uint(0), tgc.ownerCount)
	assert.Equal(t, uint(math.MaxUint16), tgc.limit)
	assert.False(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 2, called)

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
	tgc, err := New(cf, WithLimit(limit), WithInitialOwnerCount(initCount))
	assert.Nil(t, err)
	require.NotNil(t, tgc)

	// empty tgc with initialOwnerCount of 5 and limit of 10
	assert.Equal(t, uint(initCount), tgc.ownerCount)
	assert.Equal(t, uint(limit), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// add 1, with prev 5 = 6
	tgc.Add()

	assert.Equal(t, uint(initCount+1), tgc.ownerCount)
	assert.Equal(t, uint(limit), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// add 9, with prev 6 = 10 (it shouldn't go over the limit)
	for i := 0; i < 9; i += 1 {
		tgc.Add()
	}

	assert.Equal(t, uint(10), tgc.ownerCount)
	assert.Equal(t, uint(limit), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// remove 1, with prev 10 = 9
	tgc.Delete()

	assert.Equal(t, uint(9), tgc.ownerCount)
	assert.Equal(t, uint(limit), tgc.limit)
	assert.True(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 1, called)

	// remove 20, with prev 9 = 0 (it shouldn't go below 0)
	for i := 0; i < 20; i += 1 {
		tgc.Delete()
	}

	assert.Equal(t, uint(0), tgc.ownerCount)
	assert.Equal(t, uint(limit), tgc.limit)
	assert.False(t, tgc.state)
	assert.NotNil(t, tgc.onChange)
	testHelperFunctions(t, tgc)
	assert.Equal(t, 2, called)
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

func TestShutdown(t *testing.T) {
	called := 0
	cf := func(b bool) {
		called += 1
	}
	tgc, err := New(cf)
	assert.Nil(t, err)
	require.NotNil(t, tgc)

	tgc.Add()
	assert.Equal(t, 1, called)

	tgc.Add()
	tgc.Add()
	assert.Equal(t, uint(3), tgc.ownerCount)
	assert.Equal(t, 1, called)

	tgc.Shutdown()
	assert.Equal(t, uint(0), tgc.ownerCount)
	assert.Equal(t, 2, called)
}

// this functions tests if count,
func testHelperFunctions(t *testing.T, tgc *Tgc) {
	assert.Equal(t, tgc.ownerCount, tgc.Count())
	assert.Equal(t, tgc.limit, tgc.Limit())
	assert.Equal(t, tgc.state, tgc.State())
}
