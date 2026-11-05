package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Open(ctx context.Context, databaseURL string) (*gorm.DB, error) {
	dbLogger := gormlogger.NewSlogLogger(
		slog.Default().With("component", "database"),
		gormlogger.Config{
			LogLevel:                  gormlogger.Warn,
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger:               dbLogger,
		DisableAutomaticPing: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database pool: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
