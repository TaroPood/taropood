package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	checks []func(ctx context.Context) error
}

func NewHealthHandler(checks ...func(ctx context.Context) error) *HealthHandler {
	return &HealthHandler{checks: checks}
}

func (h *HealthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/healthz", h.Liveness)
	r.GET("/readyz", h.Readiness)
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	for _, check := range h.checks {
		if err := check(c.Request.Context()); err != nil {
			slog.Warn("readiness check failed", "err", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"reason": "dependency check failed",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
