package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
)

const migrationsPath = "migrations"

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load .env: %w", err)
	}

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	log, err := logger.LoadLogger(&config.Config{
		AppENV:   os.Getenv("APP_ENV"),
		LogLevel: level,
	})
	if err != nil {
		return err
	}
	slog.SetDefault(log)

	if len(os.Args) < 2 {
		return fmt.Errorf("usage: migrate <up|down|migrate|version|create|force>")
	}

	command := os.Args[1]

	// Creating SQL files doesn't require a database.
	if command == "create" {
		if len(os.Args) != 3 {
			return fmt.Errorf("usage: migrate create <name>")
		} else if err := createMigration(os.Args[2]); err != nil {
			return err
		}
		return nil
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is required in .env")
	}

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			slog.Error("failed to close migration source", "error", sourceErr)
		}
		if databaseErr != nil {
			slog.Error("failed to close migration database", "error", databaseErr)
		}
	}()

	switch command {
	case "up":
		err = m.Up()

	case "down":
		err = m.Steps(-1)

	case "migrate":
		if len(os.Args) != 3 {
			return fmt.Errorf("usage: migrate migrate <version>")
		}

		version, parseErr := strconv.ParseUint(os.Args[2], 10, 64)
		if parseErr != nil {
			return parseErr
		}
		err = m.Migrate(uint(version))

	case "version":
		version, dirty, versionErr := m.Version()
		if errors.Is(versionErr, migrate.ErrNilVersion) {
			slog.Info("No migrations applied")
			return nil
		}
		if versionErr != nil {
			return versionErr
		}
		slog.Info("migration version", "version", version, "dirty", dirty)
		return nil

	case "force":
		if len(os.Args) != 3 {
			return fmt.Errorf("usage: migrate force <version>")
		}

		version, parseErr := strconv.Atoi(os.Args[2])
		if parseErr != nil {
			return parseErr
		}
		err = m.Force(version)

	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No changes")
		return nil
	}
	if err != nil {
		return err
	}

	slog.Info("Migration completed successfully")
	return nil
}

func createMigration(name string) error {
	if name == "" || strings.ContainsAny(name, `/\`) ||
		strings.Contains(name, "..") {
		return fmt.Errorf("invalid migration name")
	}

	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return err
	}

	version := time.Now().UTC().Format("20060102150405")

	for _, direction := range []string{"up", "down"} {
		filename := fmt.Sprintf(
			"%s/%s_%s.%s.sql",
			migrationsPath, version, name, direction,
		)

		file, err := os.OpenFile(
			filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644,
		)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}

		slog.Info("migration file created", "path", filename)
	}

	return nil
}
