package userdomain

import (
	"slices"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

type OAuthProvider struct {
	valueobjects.NonEmptyString
}

var validProviders []string = []string{"google", "github"}

func NewOAuthProvider(value string) (OAuthProvider, error) {
	if !slices.Contains(validProviders, value) {
		return OAuthProvider{}, domain.ValidationError("provider", []string{"Invalid oauth provider: " + value})
	}

	nonEmpty, _ := valueobjects.NewNonEmptyString(value, "provider", "Provider", 0, 0,
		nil)

	return OAuthProvider{
		NonEmptyString: nonEmpty,
	}, nil
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
