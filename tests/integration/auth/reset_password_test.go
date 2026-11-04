package auth

import (
	"net/http"
	"strings"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/Desalutar20/lingostruct-server-go/tests/integration/testapp"
)

func TestResetPassword(t *testing.T) {
	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			app.CreateAndVerify(t, &request)

			forgotPasswordResponse := app.ForgotPassword(t, &authroutes.ForgotPasswordRequest{Email: request.Email})
			defer forgotPasswordResponse.Body.Close()

			app.AssertValidStatus(t, forgotPasswordResponse, http.StatusOK)

			resetPasswordToken := app.GetTokenFromRedis(t, testapp.ResetPassword)
			if resetPasswordToken == nil {
				t.Fatal("expected reset password token to exist in Redis")
			}

			newPassword := testapp.Password()
			response := app.ResetPassword(t, &authroutes.ResetPasswordRequest{Email: request.Email, Token: *resetPasswordToken, NewPassword: newPassword})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusOK)
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		request := testapp.ResetPasswordRequest()

		testapp.Run(t, func(app *testapp.TestApp) {
			data := []struct {
				request func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest
				field   string
				name    string
			}{
				{
					field: "email",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Email = ""
						return r
					},
					name: "empty email",
				},
				{
					field: "email",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Email = "   "
						return r
					},
					name: "whitespace email",
				},
				{
					field: "email",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Email = "invalid"
						return r
					},
					name: "invalid email",
				},
				{
					field: "email",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Email = strings.Repeat("t", valueobjects.EmailMaxLength) + "@gmail.com"
						return r
					},
					name: "email exceeds maximum length",
				},
				{
					field: "newPassword",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.NewPassword = ""
						return r
					},
					name: "empty new password",
				},
				{
					field: "newPassword",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.NewPassword = "   "
						return r
					},
					name: "whitespace new password",
				},
				{
					field: "newPassword",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.NewPassword = strings.Repeat("t", userdomain.PasswordMinLength-1)
						return r
					},
					name: "new password too short",
				},
				{
					field: "newPassword",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.NewPassword = strings.Repeat("t", userdomain.PasswordMaxLength+1)
						return r
					},
					name: "new password too long",
				},
				{
					field: "token",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Token = ""
						return r
					},
					name: "empty token",
				},
				{
					field: "token",
					request: func(r authroutes.ResetPasswordRequest) authroutes.ResetPasswordRequest {
						r.Token = "   "
						return r
					},
					name: "whitespace token",
				},
			}

			for _, data := range data {
				t.Run(data.name, func(t *testing.T) {
					t.Parallel()

					r := data.request(request)
					response := app.ResetPassword(t, &r)
					defer response.Body.Close()

					app.AssertValidStatus(t, response, http.StatusBadRequest)
					app.AssertValidationError(t, response, data.field)
				})
			}
		})

	})

	t.Run("Should Be Rate Limited", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.ResetPasswordRequest()
			request.Email = "invalid"

			for i := 0; i < app.Config.RateLimit.ResetPassword; i++ {
				response := app.ResetPassword(t, &request)
				defer response.Body.Close()

				app.AssertValidStatus(t, response, http.StatusBadRequest)
			}

			response := app.ResetPassword(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusTooManyRequests)

		})
	})
}
