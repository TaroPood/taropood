package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/interface/http/handler"
	"github.com/TaroPood/taropood/pkg/api"
)

func SetupServer(httpCfg *config.HTTPConfig, dbReady func() bool) *api.Server{
	mux := http.NewServeMux()
	healthHandler := handler.NewHealthHandler(func(ctx context.Context) error {
		if dbReady() {
			return nil
		}
		return errors.New("database not ready")
	})
	healthHandler.RegisterRoutes(mux)
	
	server := api.NewServer(httpCfg.Addr, 
		api.WithReadTimeout(httpCfg.ReadTimeout), 
		api.WithWriteTimeout(httpCfg.WriteTimeout),
		api.WithShutdownTimeout(30*time.Second),
        api.WithBaseContext(func(_ net.Listener) context.Context {
            return context.Background()
        }),
        api.WithOnShutdown(func(ctx context.Context) error {
            slog.Info("running shutdown hooks")
            return nil
        }),
	)
	server.Setup(mux)
	return server
	
}