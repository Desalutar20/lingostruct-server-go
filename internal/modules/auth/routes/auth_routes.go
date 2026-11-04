package authroutes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/dto"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
)

type Router struct {
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
	logoutCommandHandler         cqrs.CommandHandler[*authusecases.LogoutCommand, struct{}]
}

func (router *Router) SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
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
	middlewares *routes.Middlewares,
	logger *slog.Logger,
	config *config.ApplicationConfig,
	rateLimitConfig *config.RateLimitConfig,
	signUpCommandHandler cqrs.CommandHandler[*authusecases.SignUpCommand, struct{}],
	signInCommandHandler cqrs.CommandHandler[*authusecases.SignInCommand, *authdto.SignInResultDto],
	verifyAccountCommandHandler cqrs.CommandHandler[*authusecases.VerifyAccountCommand, struct{}],
	forgotPasswordCommandHandler cqrs.CommandHandler[*authusecases.ForgotPasswordCommand, struct{}],
	resetPasswordCommandHandler cqrs.CommandHandler[*authusecases.ResetPasswordCommand, struct{}],
	generateOAuthUrlQueryHandler cqrs.QueryHandler[*authusecases.GenerateOAuthUrlQuery, *authdto.OAuthSignInResultDto],
	oauthSignInCommandHandler cqrs.CommandHandler[*authusecases.OAuthSignInCommand, uuid.UUID],
	logoutCommandHandler cqrs.CommandHandler[*authusecases.LogoutCommand, struct{}],
) *Router {
	var authRouter = &Router{
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
		logoutCommandHandler,
	}

	signUpLimiter := httprate.NewRateLimiter(rateLimitConfig.SignUp, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))
	signInLimiter := httprate.NewRateLimiter(rateLimitConfig.SignIn, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))
	verifyAccountLimiter := httprate.NewRateLimiter(rateLimitConfig.VerifyAccount, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))
	forgotPasswordLimiter := httprate.NewRateLimiter(rateLimitConfig.ForgotPassword, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))
	resetPasswordLimiter := httprate.NewRateLimiter(rateLimitConfig.ResetPassword, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))
	logoutLimiter := httprate.NewRateLimiter(rateLimitConfig.Logout, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))

	router.Route("/auth", func(r chi.Router) {
		r.With(signUpLimiter.Handler).Post("/sign-up", authRouter.signUp)
		r.With(signUpLimiter.Handler).Post("/sign-up", authRouter.signUp)
		r.With(signInLimiter.Handler).Post("/sign-in", authRouter.signIn)
		r.With(verifyAccountLimiter.Handler).Post("/verify-account", authRouter.verifyAccount)
		r.With(forgotPasswordLimiter.Handler).Post("/forgot-password", authRouter.forgotPassword)
		r.With(resetPasswordLimiter.Handler).Post("/reset-password", authRouter.resetPassword)
		r.With(signInLimiter.Handler).Get("/{provider}", authRouter.getOAuthUrl)
		r.With(signInLimiter.Handler).Get("/{provider}/callback", authRouter.oauthSignIn)
		r.With(logoutLimiter.Handler).With(middlewares.Authenticate).Post("/logout", authRouter.logout)
	})

	return authRouter
}
