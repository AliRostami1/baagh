package core

import (
	"sync"
	"testing"

	"github.com/AliRostami1/baagh/pkg/tgc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/gpiod"
)

func TestRegistry(t *testing.T) {
	iReg := registry{
		chips:   map[string]map[int]*item{},
		lines:   map[string]int{},
		RWMutex: &sync.RWMutex{},
	}

	var it *item = &item{
		Line:     &gpiod.Line{},
		RWMutex:  &sync.RWMutex{},
		chipName: "",
		state:    0,
		offset:   0,
		events:   &eventRegistry{},
		options:  &ItemOptions{},
		tgc:      &tgc.Tgc{},
	}

	var chips = []string{"chip0", "chip1", "chip2"}

	// fail: get a value that doesn't exist
	i1, err := iReg.Get(chips[0], 1)
	assert.NotNil(t, err)
	require.Nil(t, i1)

	// success: create a value that doesn't exist
	err = iReg.Add(chips[0], 1, it)
	assert.Nil(t, err)

	// success: get a value that exists
	i2, err := iReg.Get(chips[0], 1)
	assert.Nil(t, err)
	require.NotNil(t, i2)
	assert.Equal(t, it, i2)

	// fail: create a value that exists
	err = iReg.Add(chips[0], 1, it)
	assert.NotNil(t, err)

	for _, chip := range chips[1:] {
		for i := 0; i < 10; i++ {
			iReg.Add(chip, i, it)
		}
	}

	var calledTotal int
	var calledChip0 int
	var calledChip1 int
	var calledChip2 int
	for _, chip := range chips {
		iReg.ForEach(chip, func(offset int, item *item) {
			assert.Equal(t, it, item)
			// assert.Equal(t, i.offset, offset)
			switch chip {
			case "chip0":
				calledChip0 += 1
			case "chip1":
				calledChip1 += 1
			case "chip2":
				calledChip2 += 1
			}
			calledTotal += 1
		})
	}
	assert.Equal(t, 1, calledChip0)
	assert.Equal(t, 10, calledChip1)
	assert.Equal(t, 10, calledChip2)
	assert.Equal(t, 21, calledTotal)

	calledTotal = 0
	calledChip0 = 0
	calledChip1 = 0
	calledChip2 = 0

	for _, chip := range chips {
		rep, err := iReg.GetAll(chip)
		assert.Nil(t, err)
		require.NotNil(t, rep)

		for _, item := range rep {
			assert.Equal(t, it, item)
			// assert.Equal(t, i.offset, offset)
			switch chip {
			case "chip0":
				calledChip0 += 1
			case "chip1":
				calledChip1 += 1
			case "chip2":
				calledChip2 += 1
			}
			calledTotal += 1
		}
	}
}
