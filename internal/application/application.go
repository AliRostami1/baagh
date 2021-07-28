package application

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/pkg/db"
	"github.com/AliRostami1/baagh/pkg/logger"
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
	// get the logger
	logger, err := logger.New()
	if err != nil {
		return nil, err
	}

	// get the config
	config, _ := config.New()
	// we temporarily ignore this check so it doesn't terminate the program, untill we add config support
	// if err != nil {
	// 	return nil, nil
	// }

	// this is the application context, it will determine when the application will exit
	ctx, cancelCtx := context.WithCancel(context.Background())

	// calling shutdown will terminate the program
	shutdown := func(reason string) {
		logger.Info(reason)
		cancelCtx()
	}

	// here we are handling terminate signals
	go signal.Handle(func(s os.Signal) {
		shutdown(fmt.Sprintf("terminating: %v signal received", s))
	})

	// Connect to and Initialize a db instnace
	db, err := db.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to db: %v", err)
	}
	logger.Info("successfully connected to db")

	return &Application{
		Log:      logger,
		Config:   config,
		Db:       db,
		Ctx:      ctx,
		Shutdown: shutdown,
	}, nil
}
