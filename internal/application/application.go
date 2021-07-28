package application

import (
	"context"
	"fmt"
	"os"

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
	Shutdown func(string)
}

func New() (*Application, error) {
	logger := logger.GetLogger()

	ctx, cancelCtx := context.WithCancel(context.Background())

	shutdown := func(reason string) {
		logger.Info(reason)
		cancelCtx()
	}

	go signal.Handle(func(s os.Signal) {
		shutdown(fmt.Sprintf("%v signal received, terminating", s))
	})

	db, err := db.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to db: %v", err)
	}

	return &Application{
		Log:      logger,
		Config:   config.GetConfig(),
		Db:       db,
		Ctx:      ctx,
		Shutdown: shutdown,
	}, nil
}
