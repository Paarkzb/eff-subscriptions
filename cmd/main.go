package main

import (
	"eff-subscriptions/internal/app"
	"eff-subscriptions/internal/config"
	"eff-subscriptions/internal/repository/postgres"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title eff-subscriptions
// @version 1.0
// @description API server for subscription application.

// @host localhost:8180
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := config.MustRead(os.Getenv("CONFIG_PATH"))

	log := setupLogger(cfg.Env)

	pgDB, err := postgres.NewPostgresDB(cfg.PostgresDBConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer pgDB.Close()

	application := app.New(log, cfg, pgDB)

	application.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
