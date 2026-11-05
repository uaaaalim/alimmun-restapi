package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := database.Open(ctx, cfg.DatabaseURL)
	cancel()

	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get database pool", "error", err)
		os.Exit(1)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			slog.Error("failed to close database", "error", err)
			os.Exit(1)
		}
	}(sqlDB)

	slog.Info("database successfully connected")

}
