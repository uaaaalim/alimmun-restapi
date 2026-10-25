package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
	"github.com/uaaaalim/alimmun-restapi/internal/middlewares"
	"github.com/uaaaalim/alimmun-restapi/internal/router"
)

// main scheme: client --[https]--> cloudflare --[https TLS origin cert]--> nginx --[http]--> api
func main() {
	bootstrap := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		bootstrap.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log, err := logger.LoadLogger(cfg)
	if err != nil {
		bootstrap.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(log)
	log.Info("logger config imported from .env")

	// database init
	db, err := database.Connect(cfg.DatabaseURL, log)
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

	// middlewares
	handler := middlewares.LoggerMiddleware(log, router.CreateRouter(db))

	// http server init
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

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
		if err := server.Close(); err != nil {
			log.Error("failed to close HTTP server", "error", err)
			return
		}
		return
	}

	log.Info("http server stopped successfully")
}
