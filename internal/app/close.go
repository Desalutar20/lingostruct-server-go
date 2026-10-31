package app

import (
	"context"
	"errors"
	"log"
)

func (a *App) Close(ctx context.Context) error {
	shutdownErr := a.server.Shutdown(ctx)
	if shutdownErr != nil {
		if err := a.server.Close(); err != nil {
			log.Printf("Failed to force close HTTP server: %v", err)
		}
	}

	if err := a.redis.Close(); err != nil {
		return errors.Join(shutdownErr, err)
	}

	if err := a.emailService.Close(); err != nil {
		return errors.Join(shutdownErr, err)
	}

	a.pool.Close()

	return shutdownErr
}
