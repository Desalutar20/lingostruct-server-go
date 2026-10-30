package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	rd "github.com/Desalutar20/lingostruct-server-go/internal/database/redis"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	server *http.Server
	pool   *pgxpool.Pool
	redis  *redis.Client
}

func (a *App) Run() error {
	log.Printf("Server starting on %s.\n", a.server.Addr)

	return a.server.ListenAndServe()
}

func (a *App) Close(ctx context.Context) error {
	shutdownErr := a.server.Shutdown(ctx)
	if shutdownErr != nil {
		if err := a.server.Close(); err != nil {
			log.Printf("Failed to force close HTTP server: %v", err)
		}
	}

	if err := a.redis.Close(); err != nil {
		return errors.Join(shutdownErr, err)
	}

	a.pool.Close()

	return shutdownErr
}

func New(ctx context.Context, config *config.Config, logger *slog.Logger) (*App, error) {
	pool, err := postgres.New(ctx, &config.Database)
	if err != nil {
		return nil, err
	}

	redis, err := rd.New(ctx, &config.Redis)
	if err != nil {
		pool.Close()
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Application.Host, config.Application.Port),
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	return &App{
		server,
		pool,
		redis,
	}, nil
}
