package authbackgroundjobs

import (
	"context"
	"log/slog"
	"time"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
)

func (j *AuthBackgroundJobs) removeExpiredSessions(ctx context.Context) {
	t := time.NewTicker(time.Minute * 10)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if _, err := j.removeExpiredSessionCommandHandler.Handle(ctx, &authusecases.RemoveExpiredSessionsCommand{}); err != nil {
				j.logger.Error(
					"Failed to remove expired sessions",
					slog.Any("error", err),
				)
			}
		}
	}
}
