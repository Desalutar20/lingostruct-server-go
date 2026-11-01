package background

import (
	"context"
	"log/slog"
	"sync"
)

type Job struct {
	Name string
	Run  func(ctx context.Context)
}

type Jobs struct {
	logger *slog.Logger
	jobs   []Job
	wg     sync.WaitGroup
}

func (j *Jobs) Start(ctx context.Context) {
	for _, job := range j.jobs {
		j.wg.Add(1)

		go func(job Job) {
			defer j.wg.Done()

			defer func() {
				if err := recover(); err != nil {
					j.logger.Error("Panic in background job",
						slog.String("job", job.Name),
						slog.Any("error", err),
					)
				}
			}()
			job.Run(ctx)

		}(job)
	}
}

func (j *Jobs) Wait() {
	j.wg.Wait()
}

func New(logger *slog.Logger, jobs ...Job) *Jobs {
	return &Jobs{
		logger: logger,
		jobs:   jobs,
	}
}
