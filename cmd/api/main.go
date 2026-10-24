package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
	"github.com/uaaaalim/alimmun-restapi/internal/router"
)

// main scheme: client --[https]--> cloudflare --[https TLS origin cert]--> nginx --[http]--> api
func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	log, err := logger.LoadLogger(cfg)
	if err != nil {
		panic(err)
	}

	log.Info("logger config imported from .env")

	// http server init
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router.CreateRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// database init
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("database connection error", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("database handle failed", "error", err)
		os.Exit(1)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Error("database close failure", "error", err)
		}
	}(sqlDB)

	log.Info("database connected successfully")

	// graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("starting http server", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "error", err)
			os.Exit(1)
		}
		return
	case <-ctx.Done():
		stop()
		log.Info("stopping http server")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown error", "error", err)
		err := server.Close()
		if err != nil {
			return
		}
		return
	}

	log.Info("http server stopped successfully")
}
