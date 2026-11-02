package authusecases

import (
	"context"

	autherrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/errors"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/passwordutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

type ResetPasswordCommand struct {
	Email       valueobjects.Email
	Token       valueobjects.NonEmptyString
	NewPassword userdomain.Password
}

type ResetPasswordCommandHandler struct {
	userRepository UserRepository
	authCache      *authservices.AuthCache
	sessionStore   *authservices.SessionStore
}

func (h *ResetPasswordCommandHandler) Handle(ctx context.Context, command *ResetPasswordCommand) (val struct{}, err error) {
	userId, err := h.authCache.GetAndDeleteUserIdByResetPasswordToken(ctx, command.Token)
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

	if user == nil || user.Email() != command.Email || !user.IsValid() {
		return struct{}{}, autherrors.ErrInvalidToken
	}

	hashedPassword, err := passwordutil.HashPassword(&command.NewPassword)
	if err != nil {
		return struct{}{}, err
	}

	if err := h.sessionStore.RemoveAll(ctx, user.Id()); err != nil {
		return struct{}{}, err
	}

	user.UpdatePassword(hashedPassword)
	if err := h.userRepository.Update(ctx, user, nil); err != nil {
		return struct{}{}, nil
	}

	return struct{}{}, nil
}
