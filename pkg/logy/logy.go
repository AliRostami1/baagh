package logy

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

type Logger interface {
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

func New(ctx context.Context, level Level, opt ...zap.Option) (*zap.SugaredLogger, error) {
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
