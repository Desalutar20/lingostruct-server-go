package app

import (
	"context"
	"errors"
)

func (a *App) Close(ctx context.Context) error {
	var shutdownErr error

	if err := a.server.Shutdown(ctx); err != nil {
		shutdownErr = errors.Join(shutdownErr, err)

		if closeErr := a.server.Close(); closeErr != nil {
			shutdownErr = errors.Join(shutdownErr, closeErr)
		}
	}

	a.jobs.Wait()

	if err := a.redis.Close(); err != nil {
		shutdownErr = errors.Join(shutdownErr, err)
	}

	if err := a.emailService.Close(); err != nil {
		shutdownErr = errors.Join(shutdownErr, err)
	}

	a.pool.Close()

	return shutdownErr
}
