package application

import (
	"context"

	"github.com/AliRostami1/baagh/internal/config"
	"github.com/AliRostami1/baagh/internal/logy"
	"github.com/AliRostami1/baagh/pkg/grace"
	"go.uber.org/zap"
)

type Application struct {
	Log    *zap.SugaredLogger
	Config *config.Config
	// DB       *database.DB
	Ctx      context.Context
	Shutdown context.CancelFunc
}

func New() (*Application, error) {
	ctx, shutdown := grace.Shutdown()

	logger, err := logy.New(ctx, logy.InfoLevel)
	if err != nil {
		return nil, err
	}

	// get the config
	config, err := config.New(&config.ConfigOptions{
		ConfigName:  "config",
		ConfigType:  "yaml",
		ConfigPaths: []string{"/etc/baagh/"},
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
