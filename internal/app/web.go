package app

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"

	"github.com/marcelofabianov/zod/config"
	"github.com/marcelofabianov/zod/internal/handler"
	"github.com/marcelofabianov/zod/pkg/web"
)

func registerWebProvider(container *dig.Container) error {
	if err := container.Provide(web.NewServer); err != nil {
		return err
	}

	if err := container.Provide(func(
		cfg *config.Config,
		logger *slog.Logger,
	) *chi.Mux {
		router := web.NewRouter(cfg, logger)

		router.Get("/", web.IndexHandler)
		router.Get("/healthz", web.HealthCheckHandler)
		router.Get("/hello", handler.Hello)

		return router
	}); err != nil {
		return err
	}
	return nil
}
