package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/TaroPood/taropood/internal/app"
	"github.com/TaroPood/taropood/internal/config"
)


func main() {
	cfgPath := ""
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		slog.Error("cannot load config", "error", err)
		os.Exit(1)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("cannot create app", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("server started", "http_addr", cfg.HTTP.Addr)
	if err := application.Run(ctx); err != nil {
		slog.Error("app run error", "error", err)
		os.Exit(1)
	}

	slog.Info("initiating graceful shutdown...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped gracefully")
}