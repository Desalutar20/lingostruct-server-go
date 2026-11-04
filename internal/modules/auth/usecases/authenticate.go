package authusecases

import (
	"context"

	autherrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/errors"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	authtypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/types"
	"github.com/google/uuid"
)

type AuthenticateCommand struct {
	SessionId uuid.UUID
}

type AuthenticateCommandHandler struct {
	sessionStore *authservices.SessionStore
}

func (h *AuthenticateCommandHandler) Handle(ctx context.Context, command *AuthenticateCommand) (*authtypes.Session, error) {
	session, err := h.sessionStore.Get(ctx, command.SessionId)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, autherrors.ErrUnauthorized
	}

	return session, nil
}
