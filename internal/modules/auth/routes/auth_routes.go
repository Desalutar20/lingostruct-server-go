package authroutes

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/dto"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthRouter struct {
	router                       *chi.Mux
	logger                       *slog.Logger
	config                       *config.ApplicationConfig
	signUpCommandHandler         cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}]
	signInCommandHandler         cqrs.CommandHandler[*authusecases.SignInCommand, *authdto.SignInResultDto]
	verifyAccountCommandHandler  cqrs.CommandHandler[*authusecases.VerifyAccountCommand, struct{}]
	forgotPasswordCommandHandler cqrs.CommandHandler[*authusecases.ForgotPasswordCommand, struct{}]
	resetPasswordCommandHandler  cqrs.CommandHandler[*authusecases.ResetPasswordCommand, struct{}]
	generateOAuthUrlQueryHandler cqrs.QueryHandler[*authusecases.GenerateOAuthUrlQuery, *authdto.OAuthSignInResultDto]
	oauthSignInCommandHandler    cqrs.CommandHandler[*authusecases.OAuthSignInCommand, uuid.UUID]
}

func (router *AuthRouter) SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   *router.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
	http.SetCookie(w, cookie)
}

func NewV1(router *chi.Mux,
	logger *slog.Logger,
	config *config.ApplicationConfig,
	signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}],
	signInCommandHandler cqrs.CommandHandler[*authusecases.SignInCommand, *authdto.SignInResultDto],
	verifyAccountCommandHandler cqrs.CommandHandler[*authusecases.VerifyAccountCommand, struct{}],
	forgotPasswordCommandHandler cqrs.CommandHandler[*authusecases.ForgotPasswordCommand, struct{}],
	resetPasswordCommandHandler cqrs.CommandHandler[*authusecases.ResetPasswordCommand, struct{}],
	generateOAuthUrlQueryHandler cqrs.QueryHandler[*authusecases.GenerateOAuthUrlQuery, *authdto.OAuthSignInResultDto],
	oauthSignInCommandHandler cqrs.CommandHandler[*authusecases.OAuthSignInCommand, uuid.UUID],

) *AuthRouter {

	var authRouter = &AuthRouter{
		router,
		logger,
		config,
		signUpCommandHandler,
		signInCommandHandler,
		verifyAccountCommandHandler,
		forgotPasswordCommandHandler,
		resetPasswordCommandHandler,
		generateOAuthUrlQueryHandler,
		oauthSignInCommandHandler,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", authRouter.signUp)
		r.Post("/sign-in", authRouter.signIn)
		r.Post("/verify-account", authRouter.verifyAccount)
		r.Post("/forgot-password", authRouter.forgotPassword)
		r.Post("/reset-password", authRouter.resetPassword)
		r.Get("/{provider}", authRouter.getOAuthUrl)
		r.Get("/{provider}/callback", authRouter.oauthSignIn)
	})

	return authRouter
}
