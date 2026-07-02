package logger


import (
	"log/slog"
	"os"
	"strings"
)

func SetupLog(format, level string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}
	
	switch strings.ToLower(format) {
		case "json":
			handler = slog.NewJSONHandler(os.Stdout, opts)
		default:
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

	slog.SetDefault(slog.New(handler))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}