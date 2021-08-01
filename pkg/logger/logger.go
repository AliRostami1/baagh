package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}, nil
}

func (l *Logger) Warningf(template string, args ...interface{}) {
	l.Warnf(template, args)
}
