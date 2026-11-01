package authusecases

import (
	"context"

	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
)

type RemoveExpiredSessionsCommand struct {
}

type RemoveExpiredSessionCommandHandler struct {
	sessionStore *authservices.SessionStore
}

func (h *RemoveExpiredSessionCommandHandler) Handle(ctx context.Context, command *RemoveExpiredSessionsCommand) (struct{}, error) {

	if err := h.sessionStore.RemoveExpired(ctx); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
