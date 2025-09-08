package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/marcelofabianov/zod/config"
	"github.com/marcelofabianov/zod/pkg/database"
	"github.com/marcelofabianov/zod/pkg/logger"
	"github.com/marcelofabianov/zod/pkg/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: application startup failed: %v", err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cfg, err := config.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log := logger.NewSlogLogger(cfg.Logger)
	log.Info("logger initialized", "level", cfg.Logger.Level)

	db, err := database.NewPostgresConnection(cfg.DB, log)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	log.Info("application starting", "env", cfg.General.Env)

	router := web.NewRouter(cfg, log)

	router.Get("/", web.IndexHandler)
	router.Get("/healthz", web.HealthCheckHandler)

	server := web.NewServer(cfg, log, router)

	go func() {
		log.Info("server is starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Info("shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", "error", err)
	}

	log.Info("server stopped")

	return nil
}
