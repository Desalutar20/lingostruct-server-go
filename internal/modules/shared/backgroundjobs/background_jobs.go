package backgroundjobs

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/background"
)

type BackgroundJobs struct {
	logger           *slog.Logger
	outboxRepository OutboxRepository
	emailSender      EmailSender
}

func (j *BackgroundJobs) Jobs() []background.Job {
	return []background.Job{
		{
			Name: "send-outbox-emails",
			Run:  j.sendOutboxEmails,
		},
	}
}

func New(
	logger *slog.Logger,
	outboxRepository OutboxRepository,
	emailSender EmailSender,
) *BackgroundJobs {
	return &BackgroundJobs{
		logger,
		outboxRepository,
		emailSender,
	}

}
