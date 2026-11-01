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

func (h *SignInCommandHandler) Handle(ctx context.Context, command *SignInCommand) (*authdto.SignInResult, error) {
	user, err := h.userRepository.GetByEmail(ctx, command.Email, nil)
	if err != nil {
		return nil, err
	}

	if user == nil || !passwordutil.VerifyPassword(&command.Password, user.HashedPassword()) || !user.IsValid() {
		return nil, autherrors.ErrInvalidCredentials
	}

	sessionIds, err := h.sessionStore.GetSessionIds(ctx, user.Id())
	if err != nil {
		return nil, err
	}

	if len(sessionIds) >= h.config.MaxSessions {
		if err := h.sessionStore.Remove(ctx, user.Id(), sessionIds[0]); err != nil {
			return nil, err
		}
	}

	session := authdto.Session{
		UserId:    user.Id(),
		Email:     user.Email(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Role:      user.Role(),
		AvatarUrl: user.AvatarUrl(),
	}

	sessionId, err := h.sessionStore.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return &authdto.SignInResult{
		UserDto:   session.ToUserDto(),
		SessionId: sessionId,
	}, nil
}
