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

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil && os.IsNotExist(err) {
		return Config{}, fmt.Errorf(".env file not found")
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

	// required envs
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DatabaseURL is required in .env")
	}

	return cfg, nil
}
