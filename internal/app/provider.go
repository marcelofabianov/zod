package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"

	"github.com/marcelofabianov/zod/config"
)

func registerProviders(cfg *config.Config, logger *slog.Logger, db *sqlx.DB) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() *config.Config { return cfg }); err != nil {
		return nil, err
	}

	if err := registerPackage(container, cfg, logger, db); err != nil {
		logger.Error("failed to register providers package", "error", err)
		return nil, err
	}

	//... register providers ...

	if err := registerWebProvider(container); err != nil {
		logger.Error("failed to register providers web", "error", err)
		return nil, err
	}

	return container, nil
}
