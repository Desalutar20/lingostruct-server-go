package adminuserusecases

import (
	"context"

	adminuserdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/dto"
	adminusertypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/types"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/google/uuid"
)

type UseCases struct {
	GetAdminUsersQueryHandler *GetAdminUsersQueryHandler
}

type AdminUserRepository interface {
	GetAll(ctx context.Context, pagination *pagination.KeysetPagination[uuid.UUID], filters *adminusertypes.AdminUsersFilters) (*pagination.KeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID], error)
}

func New(repository AdminUserRepository) *UseCases {
	return &UseCases{
		GetAdminUsersQueryHandler: &GetAdminUsersQueryHandler{
			repository,
		},
	}
}
