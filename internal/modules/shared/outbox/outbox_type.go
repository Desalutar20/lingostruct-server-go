package outbox

import (
	"slices"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type OutboxType struct {
	valueobjects.NonEmptyString
}

var validTypes []string = []string{"email"}

func NewOutboxType(value string) (OutboxType, error) {
	return valueobjects.NewNonEmptyStringBase(value, "outboxType", "Outbox type", 0, 0,
		func(value string) error {
			if !slices.Contains(validTypes, value) {
				return domain.InternalError("Invalid outbox type: " + value)
			}
			return nil
		},
		func(value valueobjects.NonEmptyString) OutboxType {
			return OutboxType{NonEmptyString: value}
		},
	)
}

func mustOutboxType(value string) OutboxType {
	role, err := NewOutboxType(value)
	if err != nil {
		panic(err)
	}

	return role
}

func OutboxTypeEmail() OutboxType {
	return mustOutboxType("email")
}
