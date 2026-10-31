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

func (id Id) MarshalBinary() ([]byte, error) {
	return uuid.UUID(id).MarshalBinary()
}

func (id *Id) UnmarshalBinary(data []byte) error {
	var value uuid.UUID

	if err := value.UnmarshalBinary(data); err != nil {
		return err
	}

	*id = Id(value)
	return nil
}
