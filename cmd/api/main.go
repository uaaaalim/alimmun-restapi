package main

import (
	"log/slog"
	"os"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
)

func main() {
	logger.LoadTempLogger()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	err = logger.LoadLogger(cfg)
	if err != nil {
		slog.Error("failed to load logger", "error", err)
		os.Exit(1)
	}

}
