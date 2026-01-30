package postgres

import (
	"context"
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg *config.DatabseConfig) (*pgxpool.Pool, error) {
	dbConfig, err := cfg.ConnectOptions()
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("Error while creating connection to the database: %w", err)
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		connPool.Close()
		return nil, fmt.Errorf("Error while acquiring connection from the database pool: %w", err)
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		connPool.Close()
		return nil, fmt.Errorf("Could not ping database: %w", err)
	}

	return connPool, nil
}
