package authroutes

import (
	"log/slog"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/go-chi/chi/v5"
)

type AuthRouter struct {
	router               *chi.Mux
	logger               *slog.Logger
	signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}]
}

func NewV1(router *chi.Mux, logger *slog.Logger, signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}]) *AuthRouter {

	var authRouter = &AuthRouter{
		router,
		logger,
		signUpCommandHandler,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", authRouter.signUp)
	})

	return authRouter
}
