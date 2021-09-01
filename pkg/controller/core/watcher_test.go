package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	t.Log("test started")
	chipName := Chips()[0]
	watcher, err := NewWatcher(chipName, 17)
	assert.Nil(t, err)
	defer watcher.Close()

	counter := 0

	if assert.NotNil(t, watcher) {
		for itemEvent := range watcher.Watch() {
			assert.Equal(t, chipName, itemEvent.Info.Name)
			assert.Equal(t, 17, itemEvent.Info.Offset)
			assert.Equal(t, true, itemEvent.Info.Used)

			counter += 1
		}
	}
}
