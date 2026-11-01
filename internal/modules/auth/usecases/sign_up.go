package authusecases

import (
	"context"
	"errors"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/passwordutil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/sleeputil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/tokenutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	usererrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SignUpCommand struct {
	Email     valueobjects.Email
	Password  userdomain.Password
	FirstName userdomain.FirstName
	LastName  userdomain.LastName
}

type SignUpCommandHandler struct {
	pool             *pgxpool.Pool
	config           *config.ApplicationConfig
	userRepository   UserRepository
	outboxRepository OutboxRepository
	emailSender      EmailSender
	authCache        *authservices.AuthCache
}

func (h *SignUpCommandHandler) Handle(ctx context.Context, command *SignUpCommand) (val struct{}, err error) {
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		return struct{}{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	user, err := h.userRepository.GetByEmail(ctx, command.Email, tx)
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
	if err := h.userRepository.Add(ctx, &newUser, tx); err != nil {
		if errors.Is(err, usererrors.ErrUserAlreadyExists) {
			return struct{}{}, nil
		}

		return struct{}{}, err
	}

	token, err := tokenutil.Generate(64)
	if err != nil {
		return struct{}{}, err
	}

	emailMessage, err := email.BuildAccountVerificationEmail(h.config, newUser.Email(), token)
	if err != nil {
		return struct{}{}, err
	}

	serialized, err := email.SerializeMessage(emailMessage)
	if err != nil {
		return struct{}{}, err
	}

	outbox := outbox.NewOutbox(outbox.OutboxTypeEmail(), serialized)
	if err := h.outboxRepository.Add(ctx, &outbox, tx); err != nil {
		return struct{}{}, nil
	}

	if err := h.authCache.StoreUserIdByVerificationToken(ctx, newUser.Id(), token); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
