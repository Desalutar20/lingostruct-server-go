package userdomain

import (
	"slices"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type UserRole struct {
	valueobjects.NonEmptyString
}

var validRoles []string = []string{"admin", "regular"}

func NewUserRole(value string) (UserRole, error) {
	return valueobjects.NewNonEmptyStringBase(value, "role", "User role", 0, 0,
		func(value string) error {
			if !slices.Contains(validRoles, value) {
				return domain.InternalError("Invalid user role: " + value)
			}

			return nil
		},
		func(value valueobjects.NonEmptyString) UserRole {
			return UserRole{NonEmptyString: value}
		},
	)
}

func mustUserRole(value string) UserRole {
	role, err := NewUserRole(value)
	if err != nil {
		panic(err)
	}

	return role
}

func UserRoleAdmin() UserRole {
	return mustUserRole("admin")
}

func UserRoleRegular() UserRole {
	return mustUserRole("regular")
}
