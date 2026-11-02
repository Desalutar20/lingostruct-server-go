package authdto

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type OAuthSignInResultDto struct {
	Url   valueobjects.Url
	State oauth.OAuthState
}
