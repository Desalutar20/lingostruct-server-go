package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/Desalutar20/lingostruct-server-go/internal/background"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	rd "github.com/Desalutar20/lingostruct-server-go/internal/database/redis"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	server       *http.Server
	listener     net.Listener
	pool         *pgxpool.Pool
	redis        *redis.Client
	emailService *email.EmailService
	jobs         *background.Jobs
	ctx          context.Context
}

func (a *App) Run() error {
	log.Printf("Server starting on %s.\n", a.listener.Addr())

	a.jobs.Start(a.ctx)

	err := a.server.Serve(a.listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (a *App) Address() string {
	return "http://" + a.listener.Addr().String()
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

	modules := setupModules(config, pool, redis, logger, services)
	jobs := setupBackgroundJobs(logger, modules)

	router := setupRouter(modules, &config.Application, logger)
	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort(
			config.Application.Host,
			strconv.Itoa(int(config.Application.Port)),
		),
	)
	if err != nil {
		pool.Close()
		redis.Close()
		return nil, fmt.Errorf("creating listener: %w", err)
	}

	server := &http.Server{
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	return &App{
		server,
		listener,
		pool,
		redis,
		services.EmailService,
		jobs,
		ctx,
	}, nil
}
