package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppENV      string
	AppPort     string
	DatabaseURL string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file not found")
	}

	cfg := Config{
		AppENV:      os.Getenv("APP_ENV"),
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}

	// automatic substitution
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.AppENV != "development" && cfg.AppENV != "production" {
		return nil, fmt.Errorf(
			"APP_ENV must be development or production, got %q",
			cfg.AppENV,
		)
	}

	// required envs
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DatabaseURL is required in .env")
	}

	return &cfg, nil
}
