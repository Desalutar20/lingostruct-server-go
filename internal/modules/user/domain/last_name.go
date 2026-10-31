package userdomain

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type LastName struct {
	valueobjects.NonEmptyString
}

const LastNameMaxLength = 50

func NewLastName(value string) (LastName, error) {
	return valueobjects.NewNonEmptyStringBase[LastName](value, "lastName", "Last name", 0, LastNameMaxLength,
		nil,
		func(value valueobjects.NonEmptyString) LastName {
			return LastName{NonEmptyString: value}
		},
	)
}

