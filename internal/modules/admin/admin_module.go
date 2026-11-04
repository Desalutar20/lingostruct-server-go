package admin

import (
	"log/slog"

	adminuser "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	adminUserModule *adminuser.Module
}

type ModuleDependencies struct {
	Pool *pgxpool.Pool
}

func (m *Module) ConfigureRoutes(router chi.Router, logger *slog.Logger, middlewares *routes.Middlewares) {
	router.With(middlewares.Authenticate,
		middlewares.Authorize([]userdomain.UserRole{userdomain.UserRoleAdmin()})).Route("/admin", func(r chi.Router) {
		m.adminUserModule.ConfigureRoutes(r, logger)
	})
}

func New(deps ModuleDependencies) *Module {

	return &Module{
		adminuser.New(adminuser.ModuleDependencies{
			Pool: deps.Pool,
		}),
	}

}
