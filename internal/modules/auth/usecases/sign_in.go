package authusecases

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/dto"
	autherrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/errors"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/passwordutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

type SignInCommand struct {
	Email    valueobjects.Email
	Password userdomain.Password
}

type SignInCommandHandler struct {
	userRepository UserRepository
	sessionStore   *authservices.SessionStore
	config         *config.ApplicationConfig
}

func (h *SignInCommandHandler) Handle(ctx context.Context, command *SignInCommand) (*authdto.SignInResultDto, error) {
	user, err := h.userRepository.GetByEmail(ctx, command.Email, nil)
	if err != nil {
		return nil, err
	}

	if user == nil || !passwordutil.VerifyPassword(&command.Password, user.HashedPassword()) || !user.IsValid() {
		return nil, autherrors.ErrInvalidCredentials
	}

	session, sessionId, err := generateSession(ctx, user, h.sessionStore, h.config, true)
	if err != nil {
		return nil, err
	}

	return &authdto.SignInResultDto{
		UserDto:   session.ToUserDto(),
		SessionId: sessionId,
	}, nil
}
