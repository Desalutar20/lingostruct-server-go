package testapp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	rd "github.com/Desalutar20/lingostruct-server-go/internal/database/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func setupPostgres(
	ctx context.Context,
	config *config.DatabaseConfig,
) (*pgxpool.Pool, func(), error) {
	dbConfig, err := pgxpool.ParseConfig(config.ConnectionStringWithoutName())
	if err != nil {
		return nil, nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, nil, err
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		connPool.Close()
		return nil, nil, fmt.Errorf("acquiring database connection: %w", err)
	}

	dbName := pgx.Identifier{config.Name}.Sanitize()

	_, err = connection.Exec(ctx, "CREATE DATABASE "+dbName)
	if err != nil {
		connection.Release()
		connPool.Close()
		return nil, nil, fmt.Errorf("creating test database: %w", err)
	}

	cleanupDatabase := func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, _ = connection.Exec(cleanupCtx, "DROP DATABASE "+dbName)

		connection.Release()
		connPool.Close()
	}

	gooseBinary, err := exec.LookPath("goose")
	if err != nil {
		cleanupDatabase()

		return nil, nil, fmt.Errorf("finding goose binary: %w", err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		cleanupDatabase()

		return nil, nil, fmt.Errorf("getting testapp source path")
	}

	migrationDir := filepath.Join(
		filepath.Dir(filename),
		"../../../migrations",
	)

	cmd := exec.CommandContext(ctx, gooseBinary, "up")
	cmd.Env = append(os.Environ(),
		"GOOSE_DRIVER=postgres",
		"GOOSE_DBSTRING="+config.ConnectionString(),
		"GOOSE_MIGRATION_DIR="+migrationDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		cleanupDatabase()

		return nil, nil, fmt.Errorf(
			"running goose migrations: %w\n%s",
			err,
			output,
		)
	}

	dbConfig, err = pgxpool.ParseConfig(config.ConnectionString())
	if err != nil {
		cleanupDatabase()
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		cleanupDatabase()
		return nil, nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		cleanupDatabase()

		return nil, nil, fmt.Errorf("connecting to test database: %w", err)
	}

	if err := seedUsers(ctx, pool); err != nil {
		pool.Close()
		cleanupDatabase()
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		cleanupDatabase()
	}

	return pool, cleanup, nil
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) error {
	users := usersSeedData()

	rows := make([][]any, 0, len(users))

	for _, user := range users {
		rows = append(rows, []any{
			user.ID,
			user.CreatedAt,
			user.UpdatedAt,
			user.Email,
			user.FirstName,
			user.LastName,
			user.HashedPassword,
			user.Role,
			user.IsBanned,
			user.IsVerified,
			user.AvatarURL,
			user.GoogleID,
			user.GithubID,
		})
	}

	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{
			"id",
			"created_at",
			"updated_at",
			"email",
			"first_name",
			"last_name",
			"hashed_password",
			"role",
			"is_banned",
			"is_verified",
			"avatar_url",
			"google_id",
			"github_id",
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("seeding users: %w", err)
	}

	return nil
}

func setupRedis(ctx context.Context, config *config.RedisConfig) (*redis.Client, func(), error) {
	client, err := rd.New(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		defer client.Close()
		var cursor uint64

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for {
			keys, nextCursor, err := client.Scan(cleanupCtx, cursor, fmt.Sprintf("%s*", config.KeyPrefix), 100).Result()
			if err != nil {
				return
			}

			if len(keys) > 0 {
				if err := client.Unlink(cleanupCtx, keys...).Err(); err != nil {
					return
				}
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}

	return client, cleanup, nil
}
