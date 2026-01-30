package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST,notEmpty" validate:"hostname|ip"`
	Port     uint   `env:"REDIS_PORT,notEmpty" validate:"port"`
	User     string `env:"REDIS_USER,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	Database uint   `env:"REDIS_DATABASE,notEmpty" validate:"min=0,max=15"`
}

func (c *RedisConfig) ConnectOptions() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Username: c.User,
		Password: c.Password,
		DB:       int(c.Database),
	}
}
