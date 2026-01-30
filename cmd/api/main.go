package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/app"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/logger"
	"github.com/joho/godotenv"
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get caller info")
	}

	if err := godotenv.Load(path.Join(filename, "../../../.env")); err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("Received signal:", sig)
		cancel()
	}()

	config := config.New()
	logger := logger.New(&config.Application)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Application.Port))
	if err != nil {
		panic(err)
	}

	app := app.New(ctx, config, logger, &listener)

	go func() {
		logger.Info(fmt.Sprintf("Listening on port :%d", config.Application.Port))

		if err := app.Run(); err != nil {
			log.Fatal(err)

		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	done := make(chan struct{})
	go func() {
		app.Close(shutdownCtx)
		listener.Close()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("Shutdown completed")
	case <-shutdownCtx.Done():
		logger.Info("Shutdown timed out")
	}
}
