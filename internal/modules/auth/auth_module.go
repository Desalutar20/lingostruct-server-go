package auth

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authbackgroundjobs "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/backgroundjobs"
	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	UseCases       *authusecases.AuthUseCases
	BackgroundJobs *authbackgroundjobs.AuthBackgroundJobs
}

type ModuleDependencies struct {
	Router            *chi.Mux
	Pool              *pgxpool.Pool
	Redis             *redis.Client
	Logger            *slog.Logger
	AppConfig         *config.ApplicationConfig
	RedisConfig       *config.RedisConfig
	UserRepository    authusecases.UserRepository
	OutboxRepository  authusecases.OutboxRepository
	EmailSender       authusecases.EmailSender
	RedisKeyGenerator authservices.RedisKeyGenerator
}

func New(deps ModuleDependencies) *Module {
	authCache := authservices.NewAuthCache(deps.Redis, deps.AppConfig, deps.RedisKeyGenerator)
	sessionStore := authservices.NewSessionStore(deps.Redis, deps.RedisKeyGenerator, deps.AppConfig, deps.RedisConfig)

	useCases := authusecases.New(deps.Pool, deps.UserRepository, deps.OutboxRepository, deps.EmailSender, authCache, sessionStore, deps.AppConfig)
	backgroundJobs := authbackgroundjobs.New(deps.Logger, useCases.RemoveExpiredSessionCommandHandler)

	authroutes.NewV1(deps.Router, deps.Logger, useCases.SignUpCommandHandler, useCases.SignInCommandHandler)

	return &Module{
		UseCases:       useCases,
		BackgroundJobs: backgroundJobs,
	}
}
