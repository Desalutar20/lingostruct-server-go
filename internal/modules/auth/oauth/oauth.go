package oauth

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type OAuthProvider interface {
	GenerateRedirectUrl(state *OAuthState) (valueobjects.Url, error)
	IsValidState(received *OAuthState, expected *OAuthState) bool
	GetUser(ctx context.Context, code valueobjects.NonEmptyString) (*OAuthUser, error)
}
