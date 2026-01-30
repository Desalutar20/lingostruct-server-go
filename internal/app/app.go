package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/clients"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	"github.com/Desalutar20/lingostruct-server-go/internal/database/redis"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user"
	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	r "github.com/redis/go-redis/v9"
)

type App struct {
	pool     *pgxpool.Pool
	redis    *r.Client
	server   *http.Server
	listener *net.Listener
	logger   *slog.Logger
}

func (a *App) Run() error {
	err := a.server.Serve(*a.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *App) Close(ctx context.Context) {
	err := a.server.Shutdown(ctx)

	if err != nil {
		a.logger.Warn(fmt.Sprintf("Error shutting down the server: %s", err))
	}

	err = a.redis.Close()
	if err != nil {
		a.logger.Warn(fmt.Sprintf("Error closing Redis connection: %s", err))
	}

	a.pool.Close()
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger, listener *net.Listener) *App {
	database, err := postgres.New(ctx, &cfg.Database)
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.New(ctx, &cfg.Redis)
	if err != nil {
		panic(err)
	}

	emailClient, err := clients.NewEmailClient(&cfg.Smtp)
	if err != nil {
		panic(err)
	}

	userModule := user.New(database)
	authModule := auth.New(&auth.Dependencies{
		Config:      &cfg.Application,
		Repository:  &userModule.Repository,
		Redis:       redisClient,
		EmailSender: emailClient,
		Logger:      logger,
	})

	middlewares := middlewares.New(authModule.Service, &cfg.Application)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RequestID)
	r.Use(httprate.Limit(100, time.Minute, httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		httputils.ErrorResponse(w, "Too many requests", http.StatusTooManyRequests)
	})))

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/auth", authModule.V1(middlewares))
		r.Mount("/users", userModule.V1())
	})

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Minute,
		Handler:      r}

	return &App{
		pool:     database,
		redis:    redisClient,
		server:   &server,
		listener: listener,
		logger:   logger,
	}
}
