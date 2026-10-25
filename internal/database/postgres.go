package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbUrl string, log *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		TranslateError: true,
		Logger:         logger.NewGORMLogger(log),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}
