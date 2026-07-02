package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/logger"
	"github.com/TaroPood/taropood/pkg/api"
)

type App struct {
	config *config.Config
	server *api.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		config: cfg,
	}
	app.server = SetupServer(&cfg.HTTP, func() bool { return true })
	logger.SetupLog(cfg.Log.Format, cfg.Log.Level)
	slog.Info(cfg.Log.Level)
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("HTTP server starting", "addr", a.config.HTTP.Addr)
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return nil
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	slog.Info("shutting down services...")

	if a.server != nil {
		if err := a.server.Shutdown(); err != nil {
			slog.Error("HTTP shutdown error", "error", err)
		}
		slog.Info("HTTP server stopped")
	}
	slog.Info("shutdown complete")
	return nil
}
