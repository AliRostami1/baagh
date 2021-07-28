package application

import (
	"context"

	"go.uber.org/zap"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/internal/logger"
	"github.com/AliRostami1/baagh/pkg/db"
	"github.com/AliRostami1/baagh/pkg/signal"
)

type Application struct {
	Log      *zap.SugaredLogger
	Config   *config.Config
	Db       *db.Db
	Ctx      context.Context
	Shutdown func()
}

func New() Application {
	ctx, cancelCtx := context.WithCancel(context.Background())
	signal.Handle(cancelCtx)

	return Application{
		Log:      logger.GetLogger(),
		Config:   config.GetConfig(),
		Db:       db.New(ctx),
		Ctx:      ctx,
		Shutdown: cancelCtx,
	}
}
