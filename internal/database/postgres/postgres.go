package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)


func connectOptions(c *config.DatabaseConfig) (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(c.ConnectionString())
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = 10
	dbConfig.MinConns = 0
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 30
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	return dbConfig, nil
}

func New(ctx context.Context, config *config.DatabaseConfig) (*pgxpool.Pool, error) {
	options, err := connectOptions(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create database config: %w", err)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, options)
	if err != nil {
		return nil, err
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		return nil, errors.New("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		return nil, errors.New("Could not ping database")
	}

	return connPool, nil

}
