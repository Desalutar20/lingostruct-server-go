package authusecases

import (
	"context"

	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

type UserRepository interface {
	GetById(ctx context.Context, id userdomain.UserId) (*userdomain.User, error)
	GetByEmail(ctx context.Context, email valueobjects.Email) (*userdomain.User, error)
	Add(ctx context.Context, user *userdomain.User) error
}

type EmailSender interface {
	Send(message email.EmailMessage) error
}

type AuthUseCases struct {
	SignUpCommandHandler *SignUpCommandHandler
}

func New(
	UserRepository UserRepository,
	EmailSender EmailSender,
	AuthCache *authservices.AuthCache,
) *AuthUseCases {

	return &AuthUseCases{
		SignUpCommandHandler: &SignUpCommandHandler{UserRepository, EmailSender, AuthCache},
	}
}
