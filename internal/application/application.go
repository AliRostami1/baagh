package application

import (
	"context"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/pkg/grace"
	"github.com/AliRostami1/baagh/pkg/logy"
)

type Application struct {
	Log    logy.Logger
	Config *config.Config
	// DB       *database.DB
	Ctx      context.Context
	Shutdown context.CancelFunc
}

func New(logLevel logy.Level) (*Application, error) {
	ctx, shutdown := grace.Shutdown()

	logger, err := logy.New(ctx, logLevel)
	if err != nil {
		return nil, err
	}

	// get the config
	config, err := config.New(&config.ConfigOptions{
		ConfigName:  "config",
		ConfigType:  "yaml",
		ConfigPaths: []string{"/etc/baagh/", "~/.config/baagh/"},
	})
	if err != nil {
		return nil, err
	}

	return &Application{
		Log:    logger,
		Config: config,
		// DB:       db,
		Ctx:      ctx,
		Shutdown: shutdown,
	}, nil
}
