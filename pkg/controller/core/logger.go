package core

import (
	"fmt"

	"github.com/AliRostami1/baagh/pkg/logy"
)

var logger logy.Logger = logy.DummyLogger{}

func SetLogger(l logy.Logger) error {
	if l == nil {
		return fmt.Errorf("logger can't be nil")
	}
	logger = l
	return nil
}
