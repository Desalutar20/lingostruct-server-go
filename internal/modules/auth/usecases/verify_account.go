package authusecases

import (
	"context"

	autherrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/errors"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type VerifyAccountCommand struct {
	Email valueobjects.Email
	Token valueobjects.NonEmptyString
}

type VerifyAccountCommandHandler struct {
	userRepository UserRepository
	authCache      *authservices.AuthCache
}

func (h *VerifyAccountCommandHandler) Handle(ctx context.Context, command *VerifyAccountCommand) (struct{}, error) {
	userId, err := h.authCache.GetAndDeleteUserIdByVerificationToken(ctx, command.Token)
	if err != nil {
		return struct{}{}, err
	}

	if userId == nil {
		return struct{}{}, autherrors.ErrInvalidToken
	}

	user, err := h.userRepository.GetById(ctx, *userId, nil)
	if err != nil {
		return struct{}{}, err
	}

	if user == nil || user.Email() != command.Email {
		return struct{}{}, autherrors.ErrInvalidToken
	}

	if err := user.Verify(); err != nil {
		return struct{}{}, err
	}

	if err := h.userRepository.Update(ctx, user, nil); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
