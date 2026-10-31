package userdomain

import (
	"slices"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type OAuthProvider struct {
	valueobjects.NonEmptyString
}

var validProviders = []string{"google", "github"}

func NewOAuthProvider(value string) (OAuthProvider, error) {
	return valueobjects.NewNonEmptyStringBase(value, "provider", "OAuth provider", 0, 0,
		func(value string) error {
			if !slices.Contains(validProviders, value) {
				return domain.ValidationError("provider", []string{"Invalid OAuth provider: " + value})
			}
			return nil
		},
		func(value valueobjects.NonEmptyString) OAuthProvider {
			return OAuthProvider{NonEmptyString: value}
		},
	)
}

func mustOAuthProvider(value string) OAuthProvider {
	provider, err := NewOAuthProvider(value)
	if err != nil {
		panic(err)
	}

	return provider
}

func OAuthProviderGoogle() OAuthProvider {
	return mustOAuthProvider("google")
}

func OAuthProviderGithub() OAuthProvider {
	return mustOAuthProvider("github")
}
