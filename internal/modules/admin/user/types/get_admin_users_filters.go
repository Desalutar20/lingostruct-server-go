package adminusertypes

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type AdminUsersFilters struct {
	Search     *valueobjects.NonEmptyString
	IsVerified *bool
	IsBanned   *bool
}
