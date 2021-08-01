package application

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/AliRostami1/baagh/pkg/config"
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
	config, err := config.New(&config.ConfigOptions{
		ConfigName:  "config",
		ConfigType:  "yaml",
		ConfigPaths: []string{"/etc/baagh/", "$HOME/.baagh", "."},
	})
	// we temporarily ignore this check so it doesn't terminate the program, untill we add config support
	if err != nil {
		return nil, err
	}

	// this is the application context, it will determine when the application will exit
	ctx, cancelCtx := context.WithCancel(context.Background())

	// calling shutdown will terminate the program
	shutdown := func(reason string) {
		logger.Info(reason)
		cancelCtx()
	}

	// here we are handling terminate signals
	signal.ShutdownHandler(shutdown)

	// Connect to and Initialize a db instnace
	db, err := db.New(ctx, config.GetString("redis_url"))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to db: %v", err)
	}

	return &Application{
		Log:      logger,
		Config:   config,
		Db:       db,
		Ctx:      ctx,
		Shutdown: shutdown,
	}, nil
}
