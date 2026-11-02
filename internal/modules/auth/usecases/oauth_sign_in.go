package authusecases

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	authservices "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/services"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

type OAuthSignInCommand struct {
	Provider      userdomain.OAuthProvider
	Code          valueobjects.NonEmptyString
	ReceivedState oauth.OAuthState
	ExpectedState oauth.OAuthState
}

type OAuthSignInCommandHandler struct {
	userRepository       UserRepository
	oauthProviderFactory *oauth.OAuthProviderFactory
	sessionStore         *authservices.SessionStore
	config               *config.ApplicationConfig
}

func (h *OAuthSignInCommandHandler) Handle(ctx context.Context, command *OAuthSignInCommand) (uuid.UUID, error) {
	provider, err := h.oauthProviderFactory.Get(command.Provider)
	if err != nil {
		return uuid.Nil, err
	}

	oauthUser, err := provider.GetUser(ctx, command.Code)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := h.userRepository.GetByEmail(ctx, oauthUser.Email, nil)
	if err != nil {
		return uuid.Nil, err
	}

	isNew := false

	if user == nil {
		newUser := userdomain.NewUser(oauthUser.Email, nil, nil, nil)
		user = &newUser
		user.Verify()
		isNew = true
	}

	if err := user.LinkOAuthProvider(command.Provider, oauthUser.ProviderId); err != nil {
		return uuid.Nil, err
	}

	if isNew {
		if err := h.userRepository.Add(ctx, user, nil); err != nil {
			return uuid.Nil, err
		}
	} else {
		if err := h.userRepository.Update(ctx, user, nil); err != nil {
			return uuid.Nil, err
		}
	}

	_, sessionId, err := generateSession(ctx, user, h.sessionStore, h.config, !isNew)
	if err != nil {
		return uuid.Nil, err
	}

	return sessionId, nil

}
