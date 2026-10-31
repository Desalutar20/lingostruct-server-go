package userdomain

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type HashedPassword struct {
	valueobjects.NonEmptyString
}

const HashedPasswordMinLength = 40
const HashedPasswordMaxLength = 200

func NewHashedPassword(value string) (HashedPassword, error) {
	val, err := valueobjects.NewNonEmptyStringBase[HashedPassword](value, "password", "Password", HashedPasswordMinLength, HashedPasswordMaxLength,
		nil,
		func(value valueobjects.NonEmptyString) HashedPassword {
			return HashedPassword{NonEmptyString: value}
		},
	)

	if err != nil {
		return val, domain.InternalError("Invalid hashed password")
	}

	return val, nil
}
