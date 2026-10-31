package auth

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	UseCases *authusecases.AuthUseCases
}

type ModuleDependencies struct {
	Router            *chi.Mux
	Pool              *pgxpool.Pool
	Redis             *redis.Client
	Logger            *slog.Logger
	AppConfig         *config.ApplicationConfig
	UserRepository    authusecases.UserRepository
	EmailSender       authusecases.EmailSender
	RedisKeyGenerator authservices.RedisKeyGenerator
}

func New(deps ModuleDependencies) *Module {
	authCache := authservices.NewAuthCache(deps.Redis, deps.AppConfig, deps.RedisKeyGenerator)
	useCases := authusecases.New(deps.UserRepository, deps.EmailSender, authCache)

	authroutes.NewV1(deps.Router, deps.Logger, useCases.SignUpCommandHandler)

	return &Module{
		UseCases: useCases,
	}
}
