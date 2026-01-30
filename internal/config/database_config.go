package config

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabseConfig struct {
	Host     string `env:"DATABASE_HOST,notEmpty" validate:"required,hostname|ip"`
	Port     uint   `env:"DATABASE_PORT,notEmpty" validate:"required,port"`
	Name     string `env:"DATABASE_NAME,notEmpty" validate:"required"`
	User     string `env:"DATABASE_USER,notEmpty" validate:"required"`
	Password string `env:"DATABASE_PASSWORD,notEmpty" validate:"required"`
	Ssl      bool   `env:"DATABASE_SSL,notEmpty"`
}

func (c *DatabseConfig) ConnectOptions() (*pgxpool.Config, error) {
	sslMode := "disable"
	if c.Ssl {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		sslMode,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 4
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = time.Minute * 30
	cfg.HealthCheckPeriod = time.Minute
	cfg.ConnConfig.ConnectTimeout = time.Second * 5

	return cfg, nil
}
