package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppENV      string
	Port        string
	DatabaseURL string
	RedisURL    string
	LogLevel    string
}

func LoadConfig(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	cfg := Config{
		AppENV:      os.Getenv("APP_ENV"),
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}

	// Auto set
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	// Necessary envs
	if cfg.AppENV != "development" && cfg.AppENV != "production" {
		return nil, fmt.Errorf("APP_ENV should be either \"development\" or \"production\"")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	return &cfg, nil
}
