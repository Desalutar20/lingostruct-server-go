package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	rd "github.com/Desalutar20/lingostruct-server-go/internal/database/redis"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	server       *http.Server
	pool         *pgxpool.Pool
	redis        *redis.Client
	emailService *email.EmailService
}

func (a *App) Run() error {
	log.Printf("Server starting on %s.\n", a.server.Addr)

	return a.server.ListenAndServe()
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

	services, err := setupService(config)
	if err != nil {
		pool.Close()
		redis.Close()
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	v1 := chi.NewRouter()
	router.Mount("/api/v1", v1)

	setupModules(config, pool, redis, logger, v1, services)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Application.Host, config.Application.Port),
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	return &App{
		server,
		pool,
		redis,
		services.EmailService,
	}, nil
}
