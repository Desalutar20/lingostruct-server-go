package app

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/background"
)

func setupBackgroundJobs(logger *slog.Logger, modules *modules) *background.Jobs {
	jobs := modules.SharedModule.BackgroundJobs.Jobs()
	jobs = append(jobs, modules.AuthModule.BackgroundJobs.Jobs()...)

	return background.New(logger, jobs...)
}
