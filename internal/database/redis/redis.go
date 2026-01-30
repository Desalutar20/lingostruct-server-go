package redis

import (
	"context"
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(cfg.ConnectOptions())
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
