package core

import (
	"fmt"
)

type Logger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

type DummyLogger struct{}

func (d DummyLogger) Errorf(string, ...interface{}) {}
func (d DummyLogger) Warnf(string, ...interface{})  {}
func (d DummyLogger) Infof(string, ...interface{})  {}
func (d DummyLogger) Debugf(string, ...interface{}) {}

var logger Logger = DummyLogger{}

func SetLogger(l Logger) error {
	if l == nil {
		return fmt.Errorf("logger can't be nil")
	}
	logger = l
	return nil
}
