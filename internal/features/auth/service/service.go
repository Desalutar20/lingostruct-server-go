package service

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetById(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (string, error)
	Update(ctx context.Context, user *model.User) error
}

type EmailSender interface {
	Send(subject, textBody, htmlBody string, to []string) error
}

type Service struct {
	config      *config.ApplicationConfig
	repository  Repository
	redis       *redis.Client
	emailSender EmailSender
}

func New(config *config.ApplicationConfig, repository Repository, redis *redis.Client, emailSender EmailSender) *Service {

	return &Service{
		config,
		repository,
		redis,
		emailSender,
	}
}
