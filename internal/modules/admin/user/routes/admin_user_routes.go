package adminuserroutes

import (
	"log/slog"

	adminuserdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/dto"
	adminuserusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Router struct {
	logger *slog.Logger

	getAdminUsersQueryHandler cqrs.QueryHandler[*adminuserusecases.GetAdminUsersQuery, *pagination.KeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID]]
}

func NewV1(router chi.Router,
	logger *slog.Logger,
	getAdminUsersQueryHandler cqrs.QueryHandler[*adminuserusecases.GetAdminUsersQuery, *pagination.KeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID]]) *Router {

	var adminUserRouter = &Router{
		logger,
		getAdminUsersQueryHandler,
	}

	router.Route("/users", func(r chi.Router) {
		r.With(httpin.NewInput(GetAdminUsersRequest{})).Get("/", adminUserRouter.GetAdminUsers)
	})

	return adminUserRouter
}
