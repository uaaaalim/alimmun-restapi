package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const migrationsPath = "migrations"

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|migrate|version|create|force>")
	}

	command := os.Args[1]

	// Creating SQL files doesn't require a database.
	if command == "create" {
		if len(os.Args) != 3 {
			log.Fatal("usage: migrate create <name>")
		}
		if err := createMigration(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		return
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required in .env")
	}

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch command {
	case "up":
		err = m.Up()

	case "down":
		err = m.Steps(-1)

	case "migrate":
		if len(os.Args) != 3 {
			log.Fatal("usage: migrate migrate <version>")
		}

		version, parseErr := strconv.ParseUint(os.Args[2], 10, 64)
		if parseErr != nil {
			log.Fatal(parseErr)
		}
		err = m.Migrate(uint(version))

	case "version":
		version, dirty, versionErr := m.Version()
		if errors.Is(versionErr, migrate.ErrNilVersion) {
			fmt.Println("No migrations applied")
			return
		}
		if versionErr != nil {
			log.Fatal(versionErr)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)
		return

	case "force":
		if len(os.Args) != 3 {
			log.Fatal("usage: migrate force <version>")
		}

		version, parseErr := strconv.Atoi(os.Args[2])
		if parseErr != nil {
			log.Fatal(parseErr)
		}
		err = m.Force(version)

	default:
		log.Fatalf("unknown command: %s", command)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No changes")
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration completed successfully")
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

		fmt.Println("Created:", filename)
	}

	return nil
}
