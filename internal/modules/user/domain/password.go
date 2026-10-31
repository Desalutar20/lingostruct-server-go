package userdomain

import (

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type Password struct {
	valueobjects.NonEmptyString
}

const PasswordMinLength = 8
const PasswordMaxLength = 60

func NewPassword(value string) (Password, error) {
	return valueobjects.NewNonEmptyStringBase[Password](value, "password", "Password", PasswordMinLength, PasswordMaxLength,
		nil,
		func(value valueobjects.NonEmptyString) Password {
			return Password{NonEmptyString: value}
		},
	)
}

