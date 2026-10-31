package userroutes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	router *chi.Mux
	logger *slog.Logger
}

func NewV1(router *chi.Mux, logger *slog.Logger) *UserRouter {

	var userRouter = &UserRouter{
		router,
		logger,
	}

	router.Route("/users", func(r chi.Router) {
	})

	return userRouter
}
