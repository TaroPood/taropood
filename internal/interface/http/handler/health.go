package handler

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"
)

type HealthHandler struct {
    checks []func(ctx context.Context) error
}

func NewHealthHandler(checks ...func(ctx context.Context) error) *HealthHandler {
    return &HealthHandler{checks: checks}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /healthz", h.Liveness)
    mux.HandleFunc("GET /readyz", h.Readiness)
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
    h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
    for _, check := range h.checks {
        if err := check(r.Context()); err != nil {
            slog.Warn("readiness check failed", "err", err)
            h.writeJSON(w, http.StatusServiceUnavailable, map[string]string{
                "status": "not_ready",
                "reason": "dependency check failed",
            })
            return
        }
    }
    h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        slog.Warn("failed to write response", "err", err)
    }
}