package app

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/admin"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type modules struct {
	SharedModule *shared.Module
	UserModule   *user.Module
	AuthModule   *auth.Module
	AdminModule  *admin.Module
}

func setupModules(config *config.Config, pool *pgxpool.Pool, redis *redis.Client, logger *slog.Logger, services *services) *modules {
	sharedModule := shared.New(shared.ModuleDependencies{Pool: pool, Logger: logger, EmailSender: services.EmailService})
	userModule := user.New(user.ModuleDependencies{Pool: pool, RateLimitConfig: &config.RateLimit})
	authModule := auth.New(auth.ModuleDependencies{
		Pool:              pool,
		Redis:             redis,
		Logger:            logger,
		AppConfig:         &config.Application,
		RedisConfig:       &config.Redis,
		OAuthConfig:       &config.OAuth,
		RateLimitConfig:   &config.RateLimit,
		UserRepository:    userModule.UserRepository,
		OutboxRepository:  outbox.NewOutboxRepository(pool),
		EmailSender:       services.EmailService,
		RedisKeyGenerator: services.RedisKeyGenerator,
	})
	adminModule := admin.New(admin.ModuleDependencies{
		Pool: pool,
	})

	return &modules{
		SharedModule: sharedModule,
		UserModule:   userModule,
		AuthModule:   authModule,
		AdminModule:  adminModule,
	}

}
