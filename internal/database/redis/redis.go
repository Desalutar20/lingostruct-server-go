package redis

import (
	"context"
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, config *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	cmd := client.Ping(ctx)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return client, nil
}
