package testapp

import (
	"net/http"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
)

func (a *TestApp) SignUp(
	t *testing.T,
	request *authroutes.SignUpRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/sign-up")
}

func (a *TestApp) VerifyAccount(
	t *testing.T,
	request *authroutes.VerifyAccountRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/verify-account")
}

func (a *TestApp) CreateAndVerify(t *testing.T, request *authroutes.SignUpRequest) {
	t.Helper()

	signUpResponse := a.SignUp(t, request)
	defer signUpResponse.Body.Close()

	a.AssertValidStatus(t, signUpResponse, http.StatusCreated)

	token := a.GetTokenFromRedis(t, AccountVerification)
	if token == nil {
		t.Fatal("expected account verification token to exist in Redis")
	}

	response := a.VerifyAccount(t, &authroutes.VerifyAccountRequest{Email: request.Email, Token: *token})
	defer response.Body.Close()

	a.AssertValidStatus(t, response, http.StatusOK)
}

func (a *TestApp) SignIn(
	t *testing.T,
	request *authroutes.SignInRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/sign-in")
}

func (a *TestApp) ForgotPassword(
	t *testing.T,
	request *authroutes.ForgotPasswordRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/forgot-password")
}

func (a *TestApp) ResetPassword(
	t *testing.T,
	request *authroutes.ResetPasswordRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/reset-password")
}
