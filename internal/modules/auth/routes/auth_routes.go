package authroutes

import (
	"log/slog"

	authdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/dto"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/go-chi/chi/v5"
)

type AuthRouter struct {
	router               *chi.Mux
	logger               *slog.Logger
	signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}]
	signInCommandHandler cqrs.CommandHandler[*authusecases.SignInCommand, *authdto.SignInResult]
}

func NewV1(router *chi.Mux, logger *slog.Logger,
	signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}],
	signInCommandHandler cqrs.CommandHandler[*authusecases.SignInCommand, *authdto.SignInResult]) *AuthRouter {

	var authRouter = &AuthRouter{
		router,
		logger,
		signUpCommandHandler,
		signInCommandHandler,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", authRouter.signUp)
		r.Post("/sign-in", authRouter.signIn)
	})

	return authRouter
}
