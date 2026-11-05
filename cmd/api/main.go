package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/cache"
	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
)

func main() {
	logger.LoadTempLogger()

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = logger.LoadLogger(cfg)
	if err != nil {
		return fmt.Errorf("load logger: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := database.Open(ctx, cfg.DatabaseURL)
	cancel()

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database pool: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	slog.Info("database successfully connected")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	redisClient, err := cache.ConnectRedis(ctx, cfg.RedisURL)
	cancel()

	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("failed to close redis", "error", err)
		}
	}()

	slog.Info("redis client has connected")
	return nil
}
