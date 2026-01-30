package setup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupDatabase(ctx context.Context, cfg *config.DatabseConfig) (*pgxpool.Pool, func()) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get caller info")
	}

	dsnWithoutDb := fmt.Sprintf(
		"postgres://%s:%s@%s:%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	dsnWithDb := fmt.Sprintf("%s/%s", dsnWithoutDb, cfg.Name)

	config, err := pgxpool.ParseConfig(dsnWithoutDb)
	if err != nil {
		panic(err)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Errorf("creating main pool: %w", err))
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		connPool.Close()
		panic(fmt.Errorf("acquiring connection: %w", err))
	}

	_, err = connection.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.Name))
	if err != nil {
		panic(fmt.Errorf("creating database: %w", err))
	}

	pool, err := pgxpool.New(ctx, dsnWithDb)
	if err != nil {
		connection.Release()
		connPool.Close()
		panic(fmt.Errorf("connecting to new database: %w", err))
	}

	cmd := exec.Command(
		"goose",
		"up",
	)

	migrationsDir := filepath.Join(filepath.Dir(filename), "../../..", "migrations")
	fmt.Printf("migrationsDir: %v\n", migrationsDir)
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOOSE_MIGRATION_DIR=%s", migrationsDir),
		"GOOSE_DRIVER=postgres",
		fmt.Sprintf("GOOSE_DBSTRING=%s", dsnWithDb),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		connection.Release()
		connPool.Close()
		panic(fmt.Errorf("failed to run migrations: %w, output: %s", err, string(out)))
	}

	return pool, func() {
		pool.Close()
		connection.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE "%s"`, cfg.Name))
		connection.Release()
		connPool.Close()
	}
}

func setupRedis(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, func()) {
	container, err := redisContainer.Run(ctx, "redis:8")
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, nat.Port("6379/tcp"))
	if err != nil {
		panic(err)
	}

	cfg.Port = uint(port.Int())

	connString, err := container.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	options, err := redis.ParseURL(connString)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(options)

	return client, func() {
		client.Close()
		container.Stop(context.Background(), nil)
	}
}
