package shared

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/backgroundjobs"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	OutboxRepository *outbox.OutboxRepository
	BackgroundJobs   *backgroundjobs.BackgroundJobs
}

type ModuleDependencies struct {
	Pool        *pgxpool.Pool
	Logger      *slog.Logger
	EmailSender backgroundjobs.EmailSender
}

func New(deps ModuleDependencies) *Module {
	repo := outbox.NewOutboxRepository(deps.Pool)

	backgroundJobs := backgroundjobs.New(deps.Logger, repo, deps.EmailSender)

	return &Module{
		OutboxRepository: repo,
		BackgroundJobs:   backgroundJobs,
	}
}
