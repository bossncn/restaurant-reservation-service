package app

import (
	"fmt"
	"github.com/bossncn/go-boilerplate/config"
	"github.com/bossncn/go-boilerplate/internal/adapters/http"
	"go.uber.org/zap"
)

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.AppEnv == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return logger, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger, err
}

func Run(cfg *config.Config) {
	logger, err := initLogger(cfg)

	if err != nil {
		fmt.Printf("error initializing logger: %v", err)
	}

	repo := http.InitRepository(cfg)
	handler := http.InitHandler(cfg)
	service := http.InitService(cfg)
	mdw := http.InitMiddleware(cfg, logger)

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Error syncing logger", zap.Error(err))
		}
	}(logger)

	// Init App Server
	server := http.NewHTTPServer(mdw, repo, handler, service)

	if err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	server.Start()
}
