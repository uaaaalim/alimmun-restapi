package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/config"
	"github.com/uaaaalim/alimmun-restapi/internal/logger"
	"github.com/uaaaalim/alimmun-restapi/internal/router"
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

	log.Info("logger config imported from .env")

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router.CreateRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("stating http server", "address", server.Addr)
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
