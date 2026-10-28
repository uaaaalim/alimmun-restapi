package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppENV  string
	AppPort string

	DatabaseURL string

	LogLevel string

	RedisURL    string
	RedisPrefix string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppENV:      os.Getenv("APP_ENV"),
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		RedisPrefix: os.Getenv("REDIS_PREFIX"),
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
			"APP_ENV must be either \"development\" or \"production\", got %q. Make sure you've set up ENV",
			cfg.AppENV,
		)
	}

	// required envs
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DatabaseURL is required in .env")
	} else if cfg.RedisPrefix == "" || cfg.RedisURL == "" {
		return nil, fmt.Errorf("RedisURL and RedisPrefix are required in .env")
	}

	return &cfg, nil
}
