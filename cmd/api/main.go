package main

import (
	"alimmun-api/internal/config"
	"alimmun-api/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	log, err := logger.LoadLogger(cfg)
	if err != nil {
		panic(err)
	}

	log.Info("Logger config imported from .env")
}
