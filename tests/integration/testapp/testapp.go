package testapp

import (
	"context"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/app"
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/logger"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TestApp struct {
	pool              *pgxpool.Pool
	redis             *redis.Client
	redisKeyGenerator *services.RedisKeyGenerator
	http              *http.Client
	Config            *config.Config
	address           string
}

func Run(t *testing.T, fn func(app *TestApp)) {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("getting testapp source path")
	}

	configPath := filepath.Join(
		filepath.Dir(filename),
		"../../../config",
		"test.yaml",
	)

	cfg, err := config.New(configPath)
	if err != nil {
		t.Fatalf("loading test config: %v", err)
	}

	cfg.Application.Port = 0
	cfg.Database.Name = "test-" + uuid.NewString()
	cfg.Redis.KeyPrefix = uuid.NewString()

	cfg.RateLimit.SignUp = 15
	cfg.RateLimit.SignIn = 15
	cfg.RateLimit.VerifyAccount = 15
	cfg.RateLimit.ForgotPassword = 15
	cfg.RateLimit.ResetPassword = 15

	logger := logger.New(&cfg.Logger)

	pool, pgCleanup, err := setupPostgres(ctx, &cfg.Database)
	if err != nil {
		cancel()
		t.Fatalf("setting up PostgreSQL: %v", err)
	}

	redis, redisCleanup, err := setupRedis(ctx, &cfg.Redis)
	if err != nil {
		cancel()
		pgCleanup()
		t.Fatalf("setting up PostgreSQL: %v", err)
	}

	log.SetOutput(io.Discard)
	application, err := app.New(ctx, &cfg, logger)
	if err != nil {
		cancel()
		pgCleanup()
		redisCleanup()
		t.Fatalf("creating application: %v", err)
	}

	http := &http.Client{
		Timeout: time.Second * 10,
	}
	redisKeyGenerator := services.NewRedisKeyGenerator(&cfg.Redis)

	testApp := &TestApp{
		pool,
		redis,
		redisKeyGenerator,
		http,
		&cfg,
		application.Address(),
	}

	t.Cleanup(func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()

		cancel()
		application.Close(closeCtx)

		pgCleanup()
		redisCleanup()
	})

	go func() {
		_ = application.Run()
	}()

	fn(testApp)
}
