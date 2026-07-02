package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/interface/http/handler"
	"github.com/TaroPood/taropood/pkg/api"
)

// TODO: replace the no-op check with actual DB/redis health checks
// when dependencies are implemented.
func SetupServer(httpCfg *config.HTTPConfig, checks ...func(ctx context.Context) error) *api.Server{
	mux := http.NewServeMux()
	healthHandler := handler.NewHealthHandler(checks...)
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