package application

import (
	"go.uber.org/zap"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/internal/logger"
)

type Application struct {
	Log    *zap.SugaredLogger
	Config config.Config
}

func NewApplication() Application {
	return Application{
		Log:    logger.GetLogger(),
		Config: config.GetConfig(),
	}
}
