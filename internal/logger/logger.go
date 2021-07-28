package logger

import (
	"go.uber.org/zap"
)

func New() (log *zap.SugaredLogger, err error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return logger.Sugar(), nil
}
