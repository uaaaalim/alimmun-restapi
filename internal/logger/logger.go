package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
)

func LoadTempLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	slog.Info("temporary logger has loaded")
}

func LoadLogger(cfg *config.Config) error {
	level := slog.LevelInfo

	switch strings.ToLower(cfg.LogLevel) {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return fmt.Errorf("log Level %q is incorrect, supported: info, warn, error, debug", cfg.LogLevel)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	return nil
}
