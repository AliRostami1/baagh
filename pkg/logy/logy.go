package logy

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

func (l *Level) String() string {
	switch *l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		panic("wrong log level")
	}
}

func (l *Level) Set(value string) error {
	switch value {
	case "debug":
		*l = DebugLevel
	case "info":
		*l = InfoLevel
	case "warn":
		*l = WarnLevel
	case "error":
		*l = ErrorLevel
	case "dpanic":
		*l = DPanicLevel
	case "panic":
		*l = PanicLevel
	case "fatal":
		*l = FatalLevel
	default:
		return fmt.Errorf("%s is not assignable to Level", value)
	}
	return nil
}

func (l *Level) Type() string {
	return "level"
}

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

type DummyLogger struct{}

func (d DummyLogger) Errorf(string, ...interface{}) {}
func (d DummyLogger) Warnf(string, ...interface{})  {}
func (d DummyLogger) Infof(string, ...interface{})  {}
func (d DummyLogger) Debugf(string, ...interface{}) {}

func New(ctx context.Context, level Level, opts ...zap.Option) (Logger, error) {
	opts = append(opts, zap.OnFatal(zapcore.WriteThenFatal))

	zapProdConfig := zap.NewProductionConfig()
	zapProdConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(level))

	logger, err := zapProdConfig.Build(opts...)
	if err != nil {
		return nil, err
	}
	zap.RedirectStdLog(logger)

	go func() {
		defer logger.Sync()
		<-ctx.Done()
	}()

	return logger.Sugar(), nil
}
