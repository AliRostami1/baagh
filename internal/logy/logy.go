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

func New(ctx context.Context, level Level, opts ...zap.Option) (*zap.SugaredLogger, error) {
	opts = append(opts, zap.OnFatal(zapcore.WriteThenFatal))

	zapProdConfig := zap.NewProductionConfig()
	zapProdConfig.Level = zap.NewAtomicLevelAt(level)

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
