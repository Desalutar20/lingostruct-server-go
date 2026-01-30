package setup

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
)

const (
	MinPasswordLength int = 8
	MaxPasswordLength int = 40

	MinFirstNameLength int = 3
	MaxFirstNameLength int = 40

	MinLastNameLength int = 3
	MaxLastNameLength int = 40
)

func (a *TestApp) SignUp(t *testing.T, body any) *http.Response {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("signUp failed: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/sign-up",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("signUp failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("signUp failed: %v", err)

	}

	return response
}

func (a *TestApp) VerifyAccount(t *testing.T, body any) *http.Response {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("verifyAccount failed: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/verify-account",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("verifyAccount failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("verifyAccount failed: %v", err)

	}

	return response
}

func (a *TestApp) SignIn(t *testing.T, body any) *http.Response {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("signIn failed: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/sign-in",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("signIn failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("signIn failed: %v", err)

	}

	return response
}

func (a *TestApp) Logout(t *testing.T, cookie *http.Cookie) *http.Response {

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/logout",
		nil,
	)
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("logout failed: %v", err)

	}

	return response
}

func (a *TestApp) ForgotPassword(t *testing.T, body any) *http.Response {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("forgotPassword failed: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/forgot-password",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("forgotPassword failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("forgotPassword failed: %v", err)

	}

	return response
}

func (a *TestApp) ResetPassword(t *testing.T, body any) *http.Response {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("signUp failed: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.addr+"/api/v1/auth/reset-password",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("resetPassword failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(req)
	if err != nil {
		t.Fatalf("resetPassword failed: %v", err)

	}

	return response
}

func (a *TestApp) CreateAndSignIn(t *testing.T, body dto.SignUpRequest) *http.Cookie {
	response := a.SignUp(t, body)
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
	}

	verificationToken := a.GetVerificationToken(t)
	if verificationToken == "" {
		t.Fatal("verification token should exists after signup")
	}

	verifyResponse := a.VerifyAccount(t, dto.VerifyAccountRequest{
		Email: body.Email,
		Token: verificationToken,
	})
	if verifyResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	signInResponse := a.SignIn(t, dto.SignInRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if signInResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	cookie := a.GetSessionCookie(t, signInResponse)

	if cookie == nil {
		t.Fatal("session cookie should exists after account verification")
	}

	return cookie
}

func (a *TestApp) GetSessionCookie(t *testing.T, response *http.Response) *http.Cookie {
	var cookie *http.Cookie

	for _, c := range response.Cookies() {
		if c.Name != a.config.SessionCookieName {
			continue
		}

		cookie = c
	}

	return cookie
}
