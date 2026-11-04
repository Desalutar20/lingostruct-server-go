package authusecases

import (
	"context"

	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

type LogoutCommand struct {
	UserId    userdomain.UserId
	SessionId uuid.UUID
}

type LogoutCommandHandler struct {
	sessionStore *authservices.SessionStore
}

func (h *LogoutCommandHandler) Handle(ctx context.Context, command *LogoutCommand) (struct{}, error) {
	if err := h.sessionStore.Remove(ctx, command.UserId, command.SessionId); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
