package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
	shutdown func(string)
}

func New(shutdown func(string)) (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return &Logger{
		SugaredLogger: logger.Sugar(),
		shutdown:      shutdown,
	}, nil
}

func (l *Logger) Warningf(template string, args ...interface{}) {
	l.Warnf(template, args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.shutdown(fmt.Sprint(args...))
	l.SugaredLogger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.shutdown(fmt.Sprintf(template, args...))
	l.SugaredLogger.Fatalf(template, args...)
}
