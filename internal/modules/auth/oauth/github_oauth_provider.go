package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
	githubEmailsURL    = "https://api.github.com/user/emails"
	githubAPIVersion   = "2026-03-10"
)

type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type githubUserResponse struct {
	ID    int64   `json:"id"`
	Email *string `json:"email"`
}

type githubUserEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type GithubOAuthProvider struct {
	httpClient *http.Client
	config     *config.OAuthProviderConfig
}

func (p *GithubOAuthProvider) GenerateRedirectUrl(state *OAuthState) (valueobjects.Url, error) {
	var u, err = url.Parse(githubAuthorizeURL)
	if err != nil {
		return valueobjects.Url{}, err
	}

	query := u.Query()
	query.Set("client_id", p.config.ClientId)
	query.Set("redirect_uri", p.config.RedirectUrl)
	query.Set("scope", "read:user user:email")
	query.Set("state", state.String())

	u.RawQuery = query.Encode()

	return valueobjects.NewUrl(u.String())
}

func (p *GithubOAuthProvider) IsValidState(received *OAuthState, expected *OAuthState) bool {
	return received == expected
}

func (p *GithubOAuthProvider) GetUser(
	ctx context.Context,
	code valueobjects.NonEmptyString,
) (*OAuthUser, error) {
	tokenRequest := struct {
		Code         string `json:"code"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURL  string `json:"redirect_uri"`
	}{
		Code:         code.Value(),
		ClientID:     p.config.ClientId,
		ClientSecret: p.config.ClientSecret,
		RedirectURL:  p.config.RedirectUrl,
	}

	tokenBody, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to encode token request: %v", err),
		)
	}

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		githubTokenURL,
		bytes.NewReader(tokenBody),
	)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to create token request: %v", err),
		)
	}

	tokenReq.Header.Set("Content-Type", "application/json")
	tokenReq.Header.Set("Accept", "application/json")

	tokenResponse, err := p.httpClient.Do(tokenReq)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: token request failed: %v", err),
		)
	}
	defer tokenResponse.Body.Close()

	if tokenResponse.StatusCode != http.StatusOK {
		return nil, domain.InternalError(
			fmt.Sprintf(
				"GithubOAuthProvider.GetUser: token request returned status %d",
				tokenResponse.StatusCode,
			),
		)
	}

	var token githubAccessTokenResponse
	if err := json.NewDecoder(tokenResponse.Body).Decode(&token); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to decode token response: %v", err),
		)
	}

	userReq, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to create user request: %v", err),
		)
	}

	userReq.Header.Set("Accept", "application/vnd.github+json")
	userReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	userReq.Header.Set("X-Github-Api-Version", githubAPIVersion)

	userResponse, err := p.httpClient.Do(userReq)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: user request failed: %v", err),
		)
	}
	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		return nil, domain.InternalError(
			fmt.Sprintf(
				"GithubOAuthProvider.GetUser: user request returned status %d",
				userResponse.StatusCode,
			),
		)
	}

	var user githubUserResponse
	if err := json.NewDecoder(userResponse.Body).Decode(&user); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to decode user response: %v", err),
		)
	}

	providerID, err := userdomain.NewProviderId(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return nil, err
	}

	if user.Email != nil {
		email, err := valueobjects.NewEmail(*user.Email)
		if err != nil {
			return nil, err
		}

		return &OAuthUser{
			ProviderId: providerID,
			Email:      email,
		}, nil
	}

	emailsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, githubEmailsURL, nil)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to create emails request: %v", err),
		)
	}

	emailsReq.Header.Set("Accept", "application/vnd.github+json")
	emailsReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	emailsReq.Header.Set("X-Github-Api-Version", githubAPIVersion)

	emailsResponse, err := p.httpClient.Do(emailsReq)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: emails request failed: %v", err),
		)
	}
	defer emailsResponse.Body.Close()

	if emailsResponse.StatusCode != http.StatusOK {
		return nil, domain.InternalError(
			fmt.Sprintf(
				"GithubOAuthProvider.GetUser: emails request returned status %d",
				emailsResponse.StatusCode,
			),
		)
	}

	var emails []githubUserEmailResponse
	if err := json.NewDecoder(emailsResponse.Body).Decode(&emails); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GithubOAuthProvider.GetUser: failed to decode emails response: %v", err),
		)
	}

	var selectedEmail *githubUserEmailResponse

	for i := range emails {
		email := &emails[i]

		if (email.Primary && email.Verified) ||
			(email.Verified && !strings.HasSuffix(email.Email, "github.com")) {
			selectedEmail = email
			break
		}
	}

	if selectedEmail == nil {
		return nil, domain.FailureError(
			"GitHub OAuth account has no verified primary email",
		)
	}

	email, err := valueobjects.NewEmail(selectedEmail.Email)
	if err != nil {
		return nil, err
	}

	return &OAuthUser{
		ProviderId: providerID,
		Email:      email,
	}, nil
}
