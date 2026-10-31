package app

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	srv "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
)

type services struct {
	EmailService      *email.EmailService
	RedisKeyGenerator *srv.RedisKeyGenerator
}

func setupService(config *config.Config) (*services, error) {
	emailService, err := email.NewEmailService(&config.Smtp)
	if err != nil {
		return nil, err
	}

	redisKeyGenerator := srv.NewRedisKeyGenerator(&config.Redis)

	return &services{
		EmailService:      emailService,
		RedisKeyGenerator: redisKeyGenerator,
	}, nil

}
