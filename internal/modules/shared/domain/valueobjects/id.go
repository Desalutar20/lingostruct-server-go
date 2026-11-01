package valueobjects

import (
	"github.com/google/uuid"
)

type Id uuid.UUID

func NewId() Id {
	return Id(uuid.New())
}

func (id Id) String() string {
	return uuid.UUID(id).String()
}
