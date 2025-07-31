package app

import (
	"context"
	"database/sql"
	"eff-subscriptions/internal/app/HTTPServer"
	"eff-subscriptions/internal/config"
	"eff-subscriptions/internal/delivery/http"
	"eff-subscriptions/internal/repository/postgres"
	"eff-subscriptions/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log        *slog.Logger
	cfg        *config.Config
	hTTPServer *HTTPServer.Server
}

func New(log *slog.Logger, cfg *config.Config, pgDB *sql.DB) *App {
	subscriptionRepository := postgres.NewSubscriptionRepository(pgDB)
	subscriptionService := service.NewSubscriptionService(log, subscriptionRepository)
	handler := http.NewHandler(log, subscriptionService)

	httpServer := HTTPServer.NewServer(cfg.HTTPConfig.Port, cfg.HTTPConfig.Timeout, handler.InitRoutes())

	return &App{
		log:        log,
		cfg:        cfg,
		hTTPServer: httpServer,
	}
}

func (app *App) MustRun() {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.log.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- app.hTTPServer.Shutdown(ctx)
	}()

	app.log.Info("starting server", "addr", app.hTTPServer.Addr(), "env", app.cfg.Env)

	app.hTTPServer.MustRun()

	err := <-shutdownError
	if err != nil {
		panic(err)
	}

	app.log.Info("server stopped")

	return
}

func (app *App) Stop(ctx context.Context) error {
	return app.hTTPServer.Shutdown(ctx)
}
