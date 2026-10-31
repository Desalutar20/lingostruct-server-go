package userdomain

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type FirstName struct {
	valueobjects.NonEmptyString
}

const FirstNameMaxLength = 30

func NewFirstName(value string) (FirstName, error) {
	return valueobjects.NewNonEmptyStringBase[FirstName](value, "firstName", "First name", 0, FirstNameMaxLength,
		nil,
		func(value valueobjects.NonEmptyString) FirstName {
			return FirstName{NonEmptyString: value}
		},
	)
}
