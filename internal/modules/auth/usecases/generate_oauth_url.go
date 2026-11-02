package authusecases

import (
	"context"

	authdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

type GenerateOAuthUrlQuery struct {
	Provider        userdomain.OAuthProvider
	AdditionalState *valueobjects.NonEmptyString
}

type GenerateOAuthUrlQueryHandler struct {
	oauthProviderFactory *oauth.OAuthProviderFactory
}

func (h *GenerateOAuthUrlQueryHandler) Handle(ctx context.Context, query *GenerateOAuthUrlQuery) (*authdto.OAuthSignInResultDto, error) {
	provider, err := h.oauthProviderFactory.Get(query.Provider)
	if err != nil {
		return nil, err
	}

	state := oauth.OAuthState{
		StateId:         uuid.New(),
		AdditionalState: query.AdditionalState,
	}

	url, err := provider.GenerateRedirectUrl(&state)
	if err != nil {
		return nil, err
	}

	return &authdto.OAuthSignInResultDto{
		Url:   url,
		State: state,
	}, nil
}
