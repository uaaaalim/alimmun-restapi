package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GORMLogger struct {
	log   *slog.Logger
	level gormlogger.LogLevel
}

func NewGORMLogger(log *slog.Logger) *GORMLogger {
	return &GORMLogger{
		log:   log.With("component", "database"),
		level: gormlogger.Warn,
	}
}

func (l *GORMLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	c := *l // copy
	c.level = level
	return &c
}

func (l *GORMLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Info {
		l.log.InfoContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *GORMLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Warn {
		l.log.WarnContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *GORMLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Error {
		l.log.ErrorContext(ctx, fmt.Sprintf(msg, args...))
	}
}

func (l *GORMLogger) Trace(
	ctx context.Context,
	begin time.Time,
	_ func() (string, int64),
	err error,
) {
	if l.level == gormlogger.Silent {
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, gorm.ErrDuplicatedKey) {
		return
	}

	duration := time.Since(begin)

	if err != nil {
		if l.level >= gormlogger.Error {
			l.log.ErrorContext(ctx, "database query failed",
				"error", err,
				"duration", duration,
			)
		}
		return
	}

	if duration > 200*time.Millisecond && l.level >= gormlogger.Warn {
		l.log.WarnContext(ctx, "slow database query",
			"duration", duration,
		)
	}
}
