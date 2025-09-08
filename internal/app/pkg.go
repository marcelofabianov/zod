package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"

	"github.com/marcelofabianov/zod/config"
	"github.com/marcelofabianov/zod/pkg/validator"
)

func registerPackage(container *dig.Container, cfg *config.Config, logger *slog.Logger, db *sqlx.DB) error {
	if err := container.Provide(func() *slog.Logger { return logger }); err != nil {
		return err
	}
	if err := container.Provide(func() *sqlx.DB { return db }); err != nil {
		return err
	}

	if err := container.Provide(validator.New); err != nil {
		return err
	}

	return nil
}
