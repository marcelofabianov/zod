package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/marcelofabianov/fault"
	"go.uber.org/dig"

	"github.com/marcelofabianov/zod/config"
)

type App struct {
	container *dig.Container
	config    *config.Config
	logger    *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger, db *sqlx.DB) (*App, error) {
	container, err := registerProviders(cfg, logger, db)
	if err != nil {
		return nil, fault.Wrap(err, "failed to register providers")
	}

	app := &App{
		container: container,
		config:    cfg,
		logger:    logger,
	}

	app.logger.Info("application container built successfully")

	return app, nil
}

func (a *App) Run() error {
	var server *http.Server

	if err := a.container.Invoke(func(srv *http.Server) {
		server = srv
	}); err != nil {
		a.logger.Error("failed to invoke http server from container", "error", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		a.logger.Info("server is starting", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("server failed to start", "error", err)
			return err
		}
	case <-stopChan:
		a.logger.Info("shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.API.WriteTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("graceful shutdown failed", "error", err)
			return err
		}
	}

	return nil
}
