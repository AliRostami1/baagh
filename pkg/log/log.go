package log

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

func New(ctx context.Context, opt ...zap.Option) (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction(opt...)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		logger.Sync()
	}()
	return logger.Sugar(), nil
}

type DummyLogger struct{}

func (d DummyLogger) Errorf(string, ...interface{}) {}
func (d DummyLogger) Warnf(string, ...interface{})  {}
func (d DummyLogger) Infof(string, ...interface{})  {}
func (d DummyLogger) Debugf(string, ...interface{}) {}
