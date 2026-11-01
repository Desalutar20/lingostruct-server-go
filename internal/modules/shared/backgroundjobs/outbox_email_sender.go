package backgroundjobs

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	"github.com/jackc/pgx/v5"
)

type EmailSender interface {
	Send(message email.EmailMessage) error
}

type OutboxRepository interface {
	GetAndLockOutboxes(ctx context.Context, outboxType outbox.OutboxType, limit uint) ([]outbox.Outbox, pgx.Tx, error)
	MarkAsProcessed(ctx context.Context, ids []outbox.OutboxId, tx pgx.Tx) error
}

func (h *BackgroundJobs) sendOutboxEmails(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	batchSize := 100
	chunkSize := 10

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			outboxes, tx, err := h.outboxRepository.GetAndLockOutboxes(ctx, outbox.OutboxTypeEmail(), uint(batchSize))
			if err != nil {
				h.logger.Error("", slog.Any("error", err))
				continue
			}

			count := len(outboxes)
			ids := make([]outbox.OutboxId, 0, count)

			for start := 0; start < count; start += chunkSize {
				end := min(start+chunkSize, count)
				chunk := outboxes[start:end]

				processedIds := h.processOutboxEmails(chunk)
				if len(processedIds) > 0 {
					ids = append(ids, processedIds...)
				}
			}

			if err := h.outboxRepository.MarkAsProcessed(ctx, ids, tx); err != nil {
				tx.Rollback(ctx)

				h.logger.Error(
					"Failed to mark outbox emails as processed",
					slog.Any("error", err),
				)

				continue
			}

			if err := tx.Commit(ctx); err != nil {
				h.logger.Error(
					"Failed to commit outbox transaction",
					slog.Any("error", err),
				)
			}
		}
	}

}

func (h *BackgroundJobs) processOutboxEmails(
	outboxes []outbox.Outbox,
) []outbox.OutboxId {
	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		ids []outbox.OutboxId
	)

	for _, item := range outboxes {
		wg.Add(1)

		go func(item outbox.Outbox) {
			defer wg.Done()

			var message email.EmailMessage

			if err := json.Unmarshal(
				[]byte(item.Data().Value()),
				&message,
			); err != nil {
				h.logger.Error(
					"Failed to decode outbox email",
					slog.Any("error", err),
					slog.Any("outbox_id", item.Id()),
				)
				return
			}

			if err := h.emailSender.Send(message); err != nil {
				h.logger.Error(
					"Failed to send outbox email",
					slog.Any("error", err),
					slog.Any("outbox_id", item.Id()),
				)
				return
			}

			mu.Lock()
			ids = append(ids, item.Id())
			mu.Unlock()
		}(item)
	}

	wg.Wait()

	return ids
}
