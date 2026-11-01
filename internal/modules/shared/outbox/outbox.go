package outbox

import (
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/google/uuid"
)

type Outbox struct {
	domain.Entity[OutboxId]
	outboxType  OutboxType
	data        valueobjects.NonEmptyString
	processedAt *time.Time
}

func (o *Outbox) OutboxType() OutboxType {
	return o.outboxType
}

func (o *Outbox) Data() valueobjects.NonEmptyString {
	return o.data
}

func (o *Outbox) ProcessedAt() *time.Time {
	return o.processedAt
}

func NewOutbox(outboxType OutboxType, data valueobjects.NonEmptyString) Outbox {
	return Outbox{
		domain.NewEntity(OutboxId(uuid.New())),
		outboxType,
		data,
		nil,
	}
}

func RestoreOutbox(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	outboxType string,
	data string,
	processedAt *time.Time,
) Outbox {
	outboxId := OutboxId(id)
	outboxT, _ := NewOutboxType(outboxType)
	outboxData, _ := valueobjects.NewNonEmptyString(data, "data", "Data", 0, 0, nil)

	return Outbox{
		Entity:      domain.RestoreEntity(outboxId, createdAt, updatedAt),
		outboxType:  outboxT,
		data:        outboxData,
		processedAt: processedAt,
	}
}
