package application

import (
	"context"
	"fmt"

	"github.com/AliRostami1/baagh/pkg/config"
	"github.com/AliRostami1/baagh/pkg/database"
	"github.com/AliRostami1/baagh/pkg/logger"
	"github.com/AliRostami1/baagh/pkg/signal"
)

type Application struct {
	Log      *logger.Logger
	Config   *config.Config
	DB       *database.DB
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
	db, err := database.New(ctx, &database.Options{
		Path:   "/var/lib/baagh/badger",
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to db: %v", err)
	}

	return &Application{
		Log:      logger,
		Config:   config,
		DB:       db,
		Ctx:      ctx,
		Shutdown: shutdown,
	}, nil
}
