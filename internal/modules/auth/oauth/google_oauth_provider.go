package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

const (
	googleAuthorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL     = "https://oauth2.googleapis.com/token"
	googleUserURL      = "https://openidconnect.googleapis.com/v1/userinfo"
)

type googleAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type googleUserResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type GoogleOAuthProvider struct {
	httpClient *http.Client
	config     *config.OAuthProviderConfig
}

func (p *GoogleOAuthProvider) GenerateRedirectUrl(state *OAuthState) (valueobjects.Url, error) {
	var u, err = url.Parse(googleAuthorizeURL)
	if err != nil {
		return valueobjects.Url{}, err
	}

	query := u.Query()
	query.Set("client_id", p.config.ClientId)
	query.Set("redirect_uri", p.config.RedirectUrl)
	query.Set("response_type", "code")
	query.Set("scope", "openid email profile")
	query.Set("access_type", "offline")
	query.Set("state", state.String())

	u.RawQuery = query.Encode()

	return valueobjects.NewUrl(u.String())
}

func (p *GoogleOAuthProvider) IsValidState(received *OAuthState, expected *OAuthState) bool {
	return received == expected
}

func (p *GoogleOAuthProvider) GetUser(
	ctx context.Context,
	code valueobjects.NonEmptyString,
) (*OAuthUser, error) {
	tokenRequest := url.Values{
		"code":          {code.String()},
		"client_id":     {p.config.ClientId},
		"client_secret": {p.config.ClientSecret},
		"redirect_uri":  {p.config.RedirectUrl},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		googleTokenURL,
		strings.NewReader(tokenRequest.Encode()),
	)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: failed to create token request: %v", err),
		)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResponse, err := p.httpClient.Do(req)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: token request failed: %v", err),
		)
	}
	defer tokenResponse.Body.Close()

	if tokenResponse.StatusCode != http.StatusOK {
		return nil, domain.InternalError(
			fmt.Sprintf(
				"GoogleOAuthProvider.GetUser: token request returned status %d",
				tokenResponse.StatusCode,
			),
		)
	}

	var token googleAccessTokenResponse
	if err := json.NewDecoder(tokenResponse.Body).Decode(&token); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: failed to decode token response: %v", err),
		)
	}

	userReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		googleUserURL,
		nil,
	)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: failed to create user request: %v", err),
		)
	}

	userReq.Header.Set("Authorization", "Bearer "+token.AccessToken)

	userResponse, err := p.httpClient.Do(userReq)
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: user request failed: %v", err),
		)
	}
	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		return nil, domain.InternalError(
			fmt.Sprintf(
				"GoogleOAuthProvider.GetUser: user request returned status %d",
				userResponse.StatusCode,
			),
		)
	}

	var user googleUserResponse
	if err := json.NewDecoder(userResponse.Body).Decode(&user); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("GoogleOAuthProvider.GetUser: failed to decode user response: %v", err),
		)
	}

	if !user.EmailVerified {
		return nil, domain.FailureError("Google user email is not verified")
	}

	providerID, err := userdomain.NewProviderId(user.Sub)
	if err != nil {
		return nil, err
	}

	email, err := valueobjects.NewEmail(user.Email)
	if err != nil {
		return nil, err
	}

	return &OAuthUser{
		ProviderId: providerID,
		Email:      email,
	}, nil
}
