package logger

import (
	"alimmun-api/internal/config"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func LoadLogger(cfg config.Config) (*slog.Logger, error) {
	var level slog.Level

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "info":
		level = slog.LevelInfo
	default:
		return nil, fmt.Errorf(
			"LOG_LEVEL must be debug, info, warn or error, got %q",
			cfg.LogLevel,
		)
	}

	options := &slog.HandlerOptions{Level: level}

	if cfg.AppENV == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, options)), nil
	}

	return slog.New(slog.NewTextHandler(os.Stdout, options)), nil
}
