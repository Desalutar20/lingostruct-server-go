package outbox

import (
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	"github.com/google/uuid"
)

type outboxDB struct {
	ID          uuid.UUID  `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	OutboxType  string     `db:"outbox_type"`
	Data        string     `db:"data"`
	ProcessedAt *time.Time `db:"processed_at"`
}

func (u outboxDB) ToDomain() *Outbox {
	return ptrutil.To(RestoreOutbox(
		u.ID,
		u.CreatedAt,
		u.UpdatedAt,
		u.OutboxType,
		u.Data,
		u.ProcessedAt,
	))
}

func fromDomain(outbox *Outbox) outboxDB {
	return outboxDB{
		uuid.UUID(outbox.Id()),
		outbox.CreatedAt(),
		outbox.UpdatedAt(),
		outbox.outboxType.Value(),
		outbox.data.Value(),
		outbox.processedAt,
	}
}
