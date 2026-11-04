package adminuserusecases

import (
	"context"

	adminuserdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/dto"
	adminusertypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/types"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/google/uuid"
)

type GetAdminUsersQuery struct {
	Pagination *pagination.KeysetPagination[uuid.UUID]
	Filters    *adminusertypes.AdminUsersFilters
}

type GetAdminUsersQueryHandler struct {
	repository AdminUserRepository
}

func (q *GetAdminUsersQueryHandler) Handle(ctx context.Context, query *GetAdminUsersQuery) (*pagination.KeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID], error) {

	return q.repository.GetAll(ctx, query.Pagination, query.Filters)
}
