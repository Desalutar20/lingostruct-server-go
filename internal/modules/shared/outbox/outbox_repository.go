package outbox

import (
	"context"
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{
		pool,
	}
}

func (r *OutboxRepository) Add(ctx context.Context, outbox *Outbox, tx pgx.Tx) error {
	var q postgres.Querier = r.pool
	if tx != nil {
		q = tx
	}

	outboxDb := fromDomain(outbox)

	_, err := q.Exec(ctx, `
		INSERT INTO outbox (
			id,
			created_at,
			updated_at,
			type,
			data
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
	`,
		outboxDb.ID,
		outboxDb.CreatedAt,
		outboxDb.UpdatedAt,
		outboxDb.OutboxType,
		outboxDb.Data,
	)

	if err != nil {
		return domain.InternalError(
			fmt.Sprintf("Failed to add outbox: %v", err),
		)
	}

	return nil
}

func (r *OutboxRepository) GetAndLockOutboxes(
	ctx context.Context,
	outboxType OutboxType,
	limit uint,
) (outboxes []Outbox, tx pgx.Tx, err error) {
	tx, err = r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT
			id,
			created_at,
			updated_at,
			processed_at,
			type,
			data
		FROM outbox
		WHERE type = $1 AND processed_at IS NULL
		ORDER BY created_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED;
	`, outboxType.Value(), limit)
	if err != nil {
		tx.Rollback(ctx)
		return nil, nil, err
	}
	defer rows.Close()

	outboxes = make([]Outbox, 0, limit)

	for rows.Next() {
		var dbOutbox outboxDB

		if err := rows.Scan(
			&dbOutbox.ID,
			&dbOutbox.CreatedAt,
			&dbOutbox.UpdatedAt,
			&dbOutbox.ProcessedAt,
			&dbOutbox.OutboxType,
			&dbOutbox.Data,
		); err != nil {
			return nil, nil, err
		}

		outboxes = append(outboxes, *dbOutbox.ToDomain())
	}

	if err := rows.Err(); err != nil {
		tx.Rollback(ctx)
		return nil, nil, err
	}

	return outboxes, tx, nil
}

func (r *OutboxRepository) MarkAsProcessed(
	ctx context.Context,
	ids []OutboxId,
	tx pgx.Tx,
) error {
	if len(ids) == 0 {
		return nil
	}

	idStrings := make([]string, 0, len(ids))

	for _, id := range ids {
		idStrings = append(idStrings, id.String())
	}

	_, err := tx.Exec(ctx, `
        UPDATE outbox
        SET updated_at = NOW(),
            processed_at = NOW()
        WHERE id = ANY($1)
    `, idStrings)

	return err
}
