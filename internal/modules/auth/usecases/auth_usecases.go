package authusecases

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetById(ctx context.Context, id userdomain.UserId, tx pgx.Tx) (*userdomain.User, error)
	GetByEmail(ctx context.Context, email valueobjects.Email, tx pgx.Tx) (*userdomain.User, error)
	Add(ctx context.Context, user *userdomain.User, tx pgx.Tx) error
}

type OutboxRepository interface {
	Add(ctx context.Context, user *outbox.Outbox, tx pgx.Tx) error
}

type EmailSender interface {
	Send(message email.EmailMessage) error
}

type AuthUseCases struct {
	SignUpCommandHandler               *SignUpCommandHandler
	SignInCommandHandler               *SignInCommandHandler
	RemoveExpiredSessionCommandHandler *RemoveExpiredSessionCommandHandler
}

func New(
	pool *pgxpool.Pool,
	userRepository UserRepository,
	outboxRepository OutboxRepository,
	emailSender EmailSender,
	authCache *authservices.AuthCache,
	sessionStore *authservices.SessionStore,
	config *config.ApplicationConfig,
) *AuthUseCases {

	return &AuthUseCases{
		SignUpCommandHandler:               &SignUpCommandHandler{pool, config, userRepository, outboxRepository, emailSender, authCache},
		SignInCommandHandler:               &SignInCommandHandler{userRepository, sessionStore, config},
		RemoveExpiredSessionCommandHandler: &RemoveExpiredSessionCommandHandler{sessionStore},
	}
}
