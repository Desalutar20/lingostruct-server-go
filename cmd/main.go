package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/app"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/logger"
)

const (
	shutdownPeriod      = 15 * time.Second
	shutdownHardPeriod  = 3 * time.Second
	readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}

		fmt.Fprintln(w, "OK")
	})

	env := os.Getenv("ENV")
	if strings.TrimSpace(env) == "" {
		env = "development"
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(
		filepath.Dir(ex),
		"../config",
		env+".yaml",
	)

	cfg, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	logger := logger.New(&cfg.Logger)
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	app, err := app.New(ongoingCtx, &cfg, logger)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-rootCtx.Done()
	stop()

	isShuttingDown.Store(true)
	logger.Info("Received shutdown signal, shutting down.")

	time.Sleep(readinessDrainDelay)
	logger.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()

	err = app.Close(shutdownCtx)
	stopOngoingGracefully()

	if err != nil {
		logger.Error("Server shutdown completed with errors", slog.Any("error", err))
	} else {
		logger.Info("Server shut down gracefully.")
	}

}
