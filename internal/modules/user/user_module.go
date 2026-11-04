package user

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userrepo "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/repository"
	userroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/routes"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	UserRepository  *userrepo.UserRepository
	rateLimitConfig *config.RateLimitConfig
}

type ModuleDependencies struct {
	Pool            *pgxpool.Pool
	RateLimitConfig *config.RateLimitConfig
}

func (m *Module) ConfigureRoutes(router *chi.Mux, logger *slog.Logger, middlewares *routes.Middlewares) {
	userroutes.NewV1(router, logger, m.rateLimitConfig, middlewares)
}

func New(deps ModuleDependencies) *Module {
	repo := userrepo.New(deps.Pool)

	return &Module{
		UserRepository:  repo,
		rateLimitConfig: deps.RateLimitConfig,
	}
}
