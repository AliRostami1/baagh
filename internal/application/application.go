package application

import (
	"go.uber.org/zap"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/internal/logger"
	"github.com/AliRostami1/baagh/pkg/db"
)

type Application struct {
	Log    *zap.SugaredLogger
	Config *config.Config
	Db     *db.Db
}

func New() Application {
	return Application{
		Log:    logger.GetLogger(),
		Config: config.GetConfig(),
		Db:     db.New(),
	}
}
