package user

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Repository repository.Repository
}

func (m *Module) V1() *chi.Mux {
	router := chi.NewRouter()

	return router
}

func New(pool *pgxpool.Pool) *Module {
	return &Module{
		Repository: *repository.New(pool),
	}
}
