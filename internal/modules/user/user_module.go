package user

import (
	"log/slog"

	userrepo "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/repository"
	userroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/routes"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	UserRepository *userrepo.UserRepository
}

type ModuleDependencies struct {
	Router *chi.Mux
	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

func New(deps ModuleDependencies) *Module {
	repo := userrepo.NewUserRepository(deps.Pool)
	userroutes.NewV1(deps.Router, deps.Logger)

	return &Module{
		UserRepository: repo,
	}
}
