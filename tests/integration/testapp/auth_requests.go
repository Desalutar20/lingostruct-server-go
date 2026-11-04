package testapp

import (
	"net/http"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

func (a *TestApp) SignUp(
	t *testing.T,
	request *authroutes.SignUpRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/sign-up", nil)
}

func (a *TestApp) VerifyAccount(
	t *testing.T,
	request *authroutes.VerifyAccountRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/verify-account", nil)
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
	return a.Request(t, request, "/api/v1/auth/sign-in", nil)
}

func (a *TestApp) ForgotPassword(
	t *testing.T,
	request *authroutes.ForgotPasswordRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/forgot-password", nil)
}

func (a *TestApp) ResetPassword(
	t *testing.T,
	request *authroutes.ResetPasswordRequest,
) *http.Response {
	t.Helper()
	return a.Request(t, request, "/api/v1/auth/reset-password", nil)
}

func (a *TestApp) CreateAndSignIn(t *testing.T, request *authroutes.SignUpRequest, role *userdomain.UserRole) []*http.Cookie {
	t.Helper()

	a.CreateAndVerify(t, request)
	if role != nil {
		_, err := a.pool.Exec(t.Context(), "UPDATE users SET role=$1 WHERE email = $2", role.Value(), request.Email)
		if err != nil {
			t.Fatal(err)
		}
	}

	response := a.SignIn(t, &authroutes.SignInRequest{Email: request.Email, Password: request.Password})
	defer response.Body.Close()

	a.AssertValidStatus(t, response, http.StatusOK)

	return response.Cookies()

}
