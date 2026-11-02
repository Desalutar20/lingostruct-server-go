package authusecases

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/tokenutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ForgotPasswordCommand struct {
	Email valueobjects.Email
}

type ForgotPasswordCommandHandler struct {
	pool             *pgxpool.Pool
	userRepository   UserRepository
	outboxRepository OutboxRepository
	authCache        *authservices.AuthCache
	config           *config.ApplicationConfig
}

func (h *ForgotPasswordCommandHandler) Handle(ctx context.Context, command *ForgotPasswordCommand) (val struct{}, err error) {
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

	if user == nil || !user.IsValid() {
		return struct{}{}, nil
	}

	token, err := tokenutil.Generate(64)
	if err != nil {
		return struct{}{}, err
	}

	if err := h.authCache.StoreUserIdByResetPasswordToken(ctx, user.Id(), token); err != nil {
		return struct{}{}, err
	}

	emailMessage, err := email.BuildResetPasswordEmail(h.config, command.Email, token)
	if err != nil {
		return struct{}{}, err
	}

	serialized, err := email.SerializeMessage(emailMessage)
	if err != nil {
		return struct{}{}, err
	}

	outbox := outbox.NewOutbox(outbox.OutboxTypeEmail(), serialized)
	if err := h.outboxRepository.Add(ctx, &outbox, tx); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}
