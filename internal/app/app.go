package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/db"
	"github.com/TaroPood/taropood/internal/logger"
	"github.com/TaroPood/taropood/pkg/api"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	server *api.Server
	db     *gorm.DB
}

func NewApp(cfg *config.Config) (*App, error) {
	logger.SetupLog(cfg.Log.Format, cfg.Log.Level)
	app := &App{
		config: cfg,
	}

	database, err := db.NewDataBase(cfg.Postgres.DSN(), cfg.Postgres.MaxOpenConns, cfg.Postgres.MaxIdleConns, cfg.Postgres.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	slog.Info("database connected", "host", cfg.Postgres.Host, "port", cfg.Postgres.Port, "db", cfg.Postgres.Db)
	app.db = database

	app.server = SetupServer(&cfg.HTTP, db.ReadinessCheck(database))
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

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				slog.Error("database close error", "error", err)
			}
			slog.Info("database connection closed")
		}
	}

	slog.Info("shutdown complete")
	return nil
}
