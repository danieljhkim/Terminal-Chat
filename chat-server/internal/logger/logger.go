package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New builds a slog.Logger using the desired level.
// Usage:
//
//	log := logger.New("debug")
//	log.Info("server started")
func New(level string) *slog.Logger {
	var lvl slog.Level

	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	case "info", "":
		lvl = slog.LevelInfo
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	return slog.New(handler)
}
