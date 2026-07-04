package api

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    "sync"
    "time"
)

var (
    ErrServerNotConfigured = errors.New("server not configured: call Setup() before Run()")
    ErrServerClosed        = errors.New("server closed")
)

type Option func(*Server)
func WithReadTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.server.ReadTimeout = d
    }
}


func WithWriteTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.server.WriteTimeout = d
    }
}


func WithIdleTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.server.IdleTimeout = d
    }
}


func WithBaseContext(fn func(net.Listener) context.Context) Option {
    return func(s *Server) {
        s.server.BaseContext = fn
    }
}


func WithShutdownTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.shutdownTimeout = d
    }
}


func WithOnError(fn func(err error)) Option {
    return func(s *Server) {
        s.onError = fn
    }
}


func WithOnShutdown(fn func(ctx context.Context) error) Option {
    return func(s *Server) {
        s.onShutdown = fn
    }
}

type Server struct {
    server          *http.Server
    shutdownTimeout time.Duration
    onError         func(err error)
    onShutdown      func(ctx context.Context) error
    configured      bool
    mu              sync.RWMutex
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        server: &http.Server{
            Addr:         addr,
            ReadTimeout:  10 * time.Second,
            WriteTimeout: 10 * time.Second,
            IdleTimeout:  60 * time.Second,
        },
        shutdownTimeout: 15 * time.Second,
        onError: func(err error) {
            slog.Error("server error", "err", err)
        },
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

func (s *Server) Setup(handler http.Handler, opts ...Option) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.server.Handler = handler

    for _, opt := range opts {
        opt(s)
    }

    s.configured = true
}

func (s *Server) Run() error {
    s.mu.RLock()
    if !s.configured {
        s.mu.RUnlock()
        return fmt.Errorf("%w: handler is nil", ErrServerNotConfigured)
    }
    s.mu.RUnlock()

    slog.Info("server starting", "addr", s.server.Addr)

    err := s.server.ListenAndServe()

    if errors.Is(err, http.ErrServerClosed) {
        slog.Info("server stopped gracefully")
        return nil
    }

    if err != nil {
        s.onError(err)
        return err
    }

    return nil
}

func (s *Server) Shutdown() error {
    s.server.SetKeepAlivesEnabled(false)

    ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
    defer cancel()

    if s.onShutdown != nil {
        if err := s.onShutdown(ctx); err != nil {
            slog.Error("shutdown hook failed", "err", err)
            return err
        }
    }


    if err := s.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server shutdown: %w", err)
    }

    return nil
}

func (s *Server) Addr() string {
    return s.server.Addr
}