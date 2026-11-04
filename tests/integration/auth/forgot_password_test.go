package auth

import (
	"net/http"
	"strings"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/tests/integration/testapp"
)

func TestForgotPassword(t *testing.T) {
	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			app.CreateAndVerify(t, &request)

			response := app.ForgotPassword(t, &authroutes.ForgotPasswordRequest{Email: request.Email})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusOK)

			resetPasswordToken := app.GetTokenFromRedis(t, testapp.ResetPassword)
			if resetPasswordToken == nil {
				t.Fatal("expected reset password token to exist in Redis")
			}
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			data := []struct {
				request func(r authroutes.ForgotPasswordRequest) authroutes.ForgotPasswordRequest
				field   string
				name    string
			}{
				{
					field: "email",
					request: func(r authroutes.ForgotPasswordRequest) authroutes.ForgotPasswordRequest {
						r.Email = ""
						return r
					},
					name: "empty email",
				},
				{
					field: "email",
					request: func(r authroutes.ForgotPasswordRequest) authroutes.ForgotPasswordRequest {
						r.Email = "   "
						return r
					},
					name: "whitespace email",
				},
				{
					field: "email",
					request: func(r authroutes.ForgotPasswordRequest) authroutes.ForgotPasswordRequest {
						r.Email = "invalid"
						return r
					},
					name: "invalid email",
				},
				{
					field: "email",
					request: func(r authroutes.ForgotPasswordRequest) authroutes.ForgotPasswordRequest {
						r.Email = strings.Repeat("t", valueobjects.EmailMaxLength) + "@gmail.com"
						return r
					},
					name: "email exceeds maximum length",
				},
			}

			for _, data := range data {
				t.Run(data.name, func(t *testing.T) {
					t.Parallel()

					r := data.request(authroutes.ForgotPasswordRequest{})
					response := app.ForgotPassword(t, &r)
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
			request := authroutes.ForgotPasswordRequest{}
			request.Email = "invalid"

			for i := 0; i < app.Config.RateLimit.ForgotPassword; i++ {
				response := app.ForgotPassword(t, &request)
				defer response.Body.Close()

				app.AssertValidStatus(t, response, http.StatusBadRequest)
			}

			response := app.ForgotPassword(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusTooManyRequests)

		})
	})
}
