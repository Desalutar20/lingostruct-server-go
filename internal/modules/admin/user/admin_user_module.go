package adminuser

import (
	"log/slog"

	adminuserrepo "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/repository"
	adminuserroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/routes"
	adminuserusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	UseCases *adminuserusecases.UseCases
}

type ModuleDependencies struct {
	Pool *pgxpool.Pool
}

func (m *Module) ConfigureRoutes(router chi.Router, logger *slog.Logger) {
	adminuserroutes.NewV1(
		router,
		logger,
		m.UseCases.GetAdminUsersQueryHandler)
}

func New(deps ModuleDependencies) *Module {
	repository := adminuserrepo.New(deps.Pool)
	useCases := adminuserusecases.New(repository)

	return &Module{
		UseCases: useCases,
	}
}
