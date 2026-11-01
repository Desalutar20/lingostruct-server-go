package authbackgroundjobs

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/background"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
)

type AuthBackgroundJobs struct {
	logger                             *slog.Logger
	removeExpiredSessionCommandHandler cqrs.CommandHandler[*authusecases.RemoveExpiredSessionsCommand, struct{}]
}

func (j *AuthBackgroundJobs) Jobs() []background.Job {
	return []background.Job{
		{
			Name: "remove-expired-sessions",
			Run:  j.removeExpiredSessions,
		},
	}
}

func New(
	logger *slog.Logger,
	removeExpiredSessionCommandHandler cqrs.CommandHandler[*authusecases.RemoveExpiredSessionsCommand, struct{}],
) *AuthBackgroundJobs {
	return &AuthBackgroundJobs{
		logger:                             logger,
		removeExpiredSessionCommandHandler: removeExpiredSessionCommandHandler,
	}

}
