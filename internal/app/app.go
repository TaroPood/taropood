package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/db"
	"github.com/TaroPood/taropood/internal/logger"
	"github.com/TaroPood/taropood/internal/repository"
	"github.com/TaroPood/taropood/internal/repository/postgres"
	"github.com/TaroPood/taropood/pkg/api"
	"gorm.io/gorm"
)

type App struct {
	config         *config.Config
	server         *api.Server
	db             *gorm.DB
	RuleRepository repository.RuleRepository
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

	app.RuleRepository = postgres.NewRuleRepository(database)

	app.server = SetupServer(&cfg.HTTP, app.RuleRepository, db.ReadinessCheck(database))
	slog.Info("application initialized", "log_level", cfg.Log.Level)

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
		if err != nil {
			return fmt.Errorf("get sql.DB for close: %w", err)
		}

		closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- sqlDB.Close()
		}()

		select {
		case err := <-done:
			if err != nil {
				return fmt.Errorf("db close: %w", err)
			}
		case <-closeCtx.Done():
			return fmt.Errorf("db close timeout")
		}

		slog.Info("database connection closed")
	}

	slog.Info("shutdown complete")
	return nil
}
