package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/zod/config"
	"github.com/marcelofabianov/zod/internal/app"
	"github.com/marcelofabianov/zod/pkg/cache"
	"github.com/marcelofabianov/zod/pkg/database"
	"github.com/marcelofabianov/zod/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: application startup failed: %v", err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return fault.Wrap(err, "failed to load .env file")
		}
	}

	cfg, err := config.LoadConfig(".")
	if err != nil {
		return fault.Wrap(err, "failed to load config")
	}

	log := logger.NewSlogLogger(cfg.Logger)
	log.Info("logger initialized", "level", cfg.Logger.Level)

	db, err := database.NewPostgresConnection(cfg.DB, log)
	if err != nil {
		return fault.Wrap(err, "failed to connect to database")
	}
	defer db.Close()

	redisClient, err := cache.NewRedisConnection(cfg.Redis, log)
	if err != nil {
		return fault.Wrap(err, "failed to connect to redis")
	}
	defer redisClient.Close()

	app, err := app.New(cfg, log, db)
	if err != nil {
		return fault.Wrap(err, "failed to initialize container application")
	}

	log.Info("application starting", "env", cfg.General.Env)

	if err := app.Run(); err != nil {
		return fault.Wrap(err, "failed to run application")
	}

	return nil
}
