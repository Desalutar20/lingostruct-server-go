package authusecases

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	authtypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/types"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/outbox"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/services/email"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetById(ctx context.Context, id userdomain.UserId, tx pgx.Tx) (*userdomain.User, error)
	GetByEmail(ctx context.Context, email valueobjects.Email, tx pgx.Tx) (*userdomain.User, error)
	Add(ctx context.Context, user *userdomain.User, tx pgx.Tx) error
	Update(ctx context.Context, user *userdomain.User, tx pgx.Tx) error
}

type OutboxRepository interface {
	Add(ctx context.Context, user *outbox.Outbox, tx pgx.Tx) error
}

type EmailSender interface {
	Send(message email.EmailMessage) error
}

type AuthUseCases struct {
	SignUpCommandHandler               *SignUpCommandHandler
	SignInCommandHandler               *SignInCommandHandler
	RemoveExpiredSessionCommandHandler *RemoveExpiredSessionCommandHandler
	VerifyAccountCommandHandler        *VerifyAccountCommandHandler
	ForgotPasswordCommandHandler       *ForgotPasswordCommandHandler
	ResetPasswordCommandHandler        *ResetPasswordCommandHandler
	GenerateOAuthUrlQueryHandler       *GenerateOAuthUrlQueryHandler
	OAuthSignInCommandHandler          *OAuthSignInCommandHandler
}

func New(
	pool *pgxpool.Pool,
	userRepository UserRepository,
	outboxRepository OutboxRepository,
	emailSender EmailSender,
	authCache *authservices.AuthCache,
	sessionStore *authservices.SessionStore,
	config *config.ApplicationConfig,
	oauthProviderFactory *oauth.OAuthProviderFactory,
) *AuthUseCases {

	return &AuthUseCases{
		SignUpCommandHandler:               &SignUpCommandHandler{pool, config, userRepository, outboxRepository, emailSender, authCache},
		SignInCommandHandler:               &SignInCommandHandler{userRepository, sessionStore, config},
		RemoveExpiredSessionCommandHandler: &RemoveExpiredSessionCommandHandler{sessionStore},
		VerifyAccountCommandHandler:        &VerifyAccountCommandHandler{userRepository, authCache},
		ForgotPasswordCommandHandler:       &ForgotPasswordCommandHandler{pool, userRepository, outboxRepository, authCache, config},
		ResetPasswordCommandHandler:        &ResetPasswordCommandHandler{userRepository, authCache, sessionStore},
		GenerateOAuthUrlQueryHandler:       &GenerateOAuthUrlQueryHandler{oauthProviderFactory},
		OAuthSignInCommandHandler:          &OAuthSignInCommandHandler{userRepository, oauthProviderFactory, sessionStore, config},
	}
}

func generateSession(ctx context.Context, user *userdomain.User, sessionStore *authservices.SessionStore, config *config.ApplicationConfig, checkSessions bool) (*authtypes.Session, uuid.UUID, error) {

	if checkSessions {
		sessionIds, err := sessionStore.GetSessionIds(ctx, user.Id())
		if err != nil {
			return nil, uuid.Nil, err
		}

		if len(sessionIds) >= config.MaxSessions {
			if err := sessionStore.Remove(ctx, user.Id(), sessionIds[0]); err != nil {
				return nil, uuid.Nil, err
			}
		}
	}

	session := authtypes.Session{
		UserId:    user.Id(),
		Email:     user.Email(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Role:      user.Role(),
		AvatarUrl: user.AvatarUrl(),
	}

	sessionId, err := sessionStore.Save(ctx, session)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return &session, sessionId, nil
}
