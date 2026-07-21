package app

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/TaroPood/taropood/internal/config"
	"github.com/TaroPood/taropood/internal/handler"
	"github.com/TaroPood/taropood/internal/middleware"
	"github.com/TaroPood/taropood/internal/repository"
	ruleuc "github.com/TaroPood/taropood/internal/usecases/rule"
	"github.com/TaroPood/taropood/pkg/api"
	"github.com/gin-gonic/gin"
)

func SetupServer(httpCfg *config.HTTPConfig, ruleRepo repository.RuleRepository, checks ...func(ctx context.Context) error) *api.Server {
	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())

	healthHandler := handler.NewHealthHandler(checks...)
	healthHandler.RegisterRoutes(&r.RouterGroup)

	ruleUC := ruleuc.NewUseCase(ruleRepo)
	ruleHandler := handler.NewRuleHandler(ruleUC)
	apiGroup := r.Group("/api")
	ruleHandler.RegisterRoutes(apiGroup)

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
	server.Setup(r)
	return server
}
