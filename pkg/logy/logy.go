package logy

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

func New(ctx context.Context, level zapcore.Level, opt ...zap.Option) (*zap.SugaredLogger, error) {
	zapProdConfig := zap.NewProductionConfig()
	zapProdConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapProdConfig.Build(opt...)
	if err != nil {
		return nil, err
	}
	zap.RedirectStdLog(logger)

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
