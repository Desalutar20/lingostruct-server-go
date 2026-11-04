package auth

import (
	"net/http"
	"strings"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/tests/integration/testapp"
)

func TestVerifyAccount(t *testing.T) {
	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()

			signUpResponse := app.SignUp(t, &request)
			defer signUpResponse.Body.Close()

			app.AssertValidStatus(t, signUpResponse, http.StatusCreated)

			token := app.GetTokenFromRedis(t, testapp.AccountVerification)
			if token == nil {
				t.Fatal("expected account verification token to exist in Redis")
			}

			response := app.VerifyAccount(t, &authroutes.VerifyAccountRequest{Email: request.Email, Token: *token})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusOK)

			user := app.GetUserByEmail(t, request.Email)
			if user == nil || !user.IsVerified {
				t.Fatal("expected user to be created and verified")
			}
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := authroutes.VerifyAccountRequest{
				Email: "test@gmail.com",
				Token: "token",
			}

			data := []struct {
				request func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest
				field   string
				name    string
			}{
				{
					field: "email",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
						r.Email = ""
						return r
					},
					name: "empty email",
				},
				{
					field: "email",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
						r.Email = "   "
						return r
					},
					name: "whitespace email",
				},
				{
					field: "email",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
						r.Email = "invalid"
						return r
					},
					name: "invalid email",
				},
				{
					field: "email",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
						r.Email = strings.Repeat("t", valueobjects.EmailMaxLength) + "@gmail.com"
						return r
					},
					name: "email exceeds maximum length",
				},
				{
					field: "token",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
						r.Token = ""
						return r
					},
					name: "empty token",
				},
				{
					field: "token",
					request: func(r authroutes.VerifyAccountRequest) authroutes.VerifyAccountRequest {
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
					response := app.VerifyAccount(t, &r)
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
			request := authroutes.VerifyAccountRequest{
				Email: "invalid",
				Token: "token",
			}

			for i := 0; i < app.Config.RateLimit.VerifyAccount; i++ {
				response := app.VerifyAccount(t, &request)
				defer response.Body.Close()

				app.AssertValidStatus(t, response, http.StatusBadRequest)
			}

			response := app.VerifyAccount(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusTooManyRequests)

		})
	})
}
