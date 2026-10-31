package authusecases

import (
	"context"
	"errors"
	"time"

	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/passwordutil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/sleeputil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/tokenutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	usererrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain/errors"
)

type SignUpCommand struct {
	Email     valueobjects.Email
	Password  userdomain.Password
	FirstName userdomain.FirstName
	LastName  userdomain.LastName
}

type SignUpCommandHandler struct {
	userRepository UserRepository
	emailSender    EmailSender
	authCache      *authservices.AuthCache
}

func (h *SignUpCommandHandler) Handle(ctx context.Context, command *SignUpCommand) (struct{}, error) {
	user, err := h.userRepository.GetByEmail(ctx, command.Email)
	if err != nil {
		return struct{}{}, err
	}

	if user != nil {
		if err := sleeputil.Sleep(ctx, time.Second*2); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	}

	hashedPassword, err := passwordutil.HashPassword(&command.Password)
	if err != nil {
		return struct{}{}, err
	}

	newUser := userdomain.NewUser(command.Email, &hashedPassword, &command.FirstName, &command.LastName)
	if err := h.userRepository.Add(ctx, &newUser); err != nil && !errors.Is(err, usererrors.UserAlreadyExists) {
		return struct{}{}, err
	}

	token, err := tokenutil.Generate(64)
	if err != nil {
		return struct{}{}, err
	}

	if err := h.authCache.StoreUserIdByVerificationToken(ctx, newUser.Id(), token); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
