package auth

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authbackgroundjobs "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/backgroundjobs"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	UseCases        *authusecases.UseCases
	BackgroundJobs  *authbackgroundjobs.AuthBackgroundJobs
	logger          *slog.Logger
	appConfig       *config.ApplicationConfig
	rateLimitConfig *config.RateLimitConfig
}

type ModuleDependencies struct {
	Pool              *pgxpool.Pool
	Redis             *redis.Client
	Logger            *slog.Logger
	AppConfig         *config.ApplicationConfig
	RedisConfig       *config.RedisConfig
	OAuthConfig       *config.OAuthConfig
	RateLimitConfig   *config.RateLimitConfig
	UserRepository    authusecases.UserRepository
	OutboxRepository  authusecases.OutboxRepository
	EmailSender       authusecases.EmailSender
	RedisKeyGenerator authservices.RedisKeyGenerator
}

func (m *Module) ConfigureRoutes(router *chi.Mux, middlewares *routes.Middlewares) {
	authroutes.NewV1(
		router,
		middlewares,
		m.logger,
		m.appConfig,
		m.rateLimitConfig,
		m.UseCases.SignUpCommandHandler,
		m.UseCases.SignInCommandHandler,
		m.UseCases.VerifyAccountCommandHandler,
		m.UseCases.ForgotPasswordCommandHandler,
		m.UseCases.ResetPasswordCommandHandler,
		m.UseCases.GenerateOAuthUrlQueryHandler,
		m.UseCases.OAuthSignInCommandHandler,
		m.UseCases.LogoutCommandHandler)
}

func New(deps ModuleDependencies) *Module {
	authCache := authservices.NewAuthCache(deps.Redis, deps.AppConfig, deps.RedisKeyGenerator)
	sessionStore := authservices.NewSessionStore(deps.Redis, deps.RedisKeyGenerator, deps.AppConfig, deps.RedisConfig)
	oauthProviderFactory := oauth.NewOAuthProviderFactory(deps.OAuthConfig)

	useCases := authusecases.New(deps.Pool, deps.UserRepository, deps.OutboxRepository, deps.EmailSender, authCache, sessionStore, deps.AppConfig, oauthProviderFactory)
	backgroundJobs := authbackgroundjobs.New(deps.Logger, useCases.RemoveExpiredSessionCommandHandler)

	return &Module{
		UseCases:        useCases,
		BackgroundJobs:  backgroundJobs,
		logger:          deps.Logger,
		appConfig:       deps.AppConfig,
		rateLimitConfig: deps.RateLimitConfig,
	}
}
