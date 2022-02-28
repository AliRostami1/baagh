package logy_test

import (
	"testing"

	"github.com/AliRostami1/baagh/pkg/logy"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	var l logy.Level

	// levels in string
	levelsStr := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
	// their corresponding level constant (undelying integer)
	levelsInt := []logy.Level{logy.DebugLevel, logy.InfoLevel, logy.WarnLevel, logy.ErrorLevel, logy.DPanicLevel, logy.PanicLevel, logy.FatalLevel}

	for i, levelStr := range levelsStr {
		err := l.Set(levelStr)
		assert.Nil(t, err)
		// it's String() method should return what we
		// passed to Set() method
		assert.Equal(t, levelStr, l.String())
		// it's actual underlying integer should match
		// it's corresponding value in levelsInt array
		assert.Equal(t, levelsInt[i], l)
	}

	assert.Equal(t, "level", l.Type())
}
