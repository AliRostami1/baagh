package logger

import (
	"log"

	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	return logger.Sugar()
}
