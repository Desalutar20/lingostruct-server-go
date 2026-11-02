package oauth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

type OAuthProviderFactory struct {
	providers map[userdomain.OAuthProvider]OAuthProvider
}

func (f *OAuthProviderFactory) Get(provider userdomain.OAuthProvider) (OAuthProvider, error) {
	p, ok := f.providers[provider]
	if !ok {
		return nil, domain.InternalError(
			fmt.Sprintf("OAuthProviderFactory.Get: unsupported provider: %s", provider),
		)
	}

	return p, nil
}
func NewOAuthProviderFactory(config *config.OAuthConfig) *OAuthProviderFactory {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &OAuthProviderFactory{
		providers: map[userdomain.OAuthProvider]OAuthProvider{
			userdomain.OAuthProviderGoogle(): &GoogleOAuthProvider{httpClient, &config.Google},
			userdomain.OAuthProviderGithub(): &GithubOAuthProvider{httpClient, &config.Github},
		},
	}
}
