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

func TestSignIn(t *testing.T) {
	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			app.CreateAndVerify(t, &request)

			response := app.SignIn(t, &authroutes.SignInRequest{Email: request.Email, Password: request.Password})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusOK)
			cookieSessionID := app.GetSessionId(t, response)
			if cookieSessionID == nil {
				t.Fatal("expected session cookie")
			}

			redisSessionID := app.GetTokenFromRedis(t, testapp.Session)
			if redisSessionID == nil {
				t.Fatal("expected session to exist in Redis")
			}

			if cookieSessionID.String() != *redisSessionID {
				t.Fatalf(
					"expected session ID %s, got %s",
					*redisSessionID,
					cookieSessionID,
				)
			}

		})
	})

	t.Run("Should Return BadRequest When User Is Not Verified", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()

			signUpResponse := app.SignUp(t, &request)
			defer signUpResponse.Body.Close()

			app.AssertValidStatus(t, signUpResponse, http.StatusCreated)

			response := app.SignIn(t, &authroutes.SignInRequest{Email: request.Email, Password: request.Password})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusBadRequest)
		})
	})

	t.Run("Should Return BadRequest When User Is Banned", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			app.CreateAndVerify(t, &request)

			app.BanUserByEmail(t, request.Email)

			response := app.SignIn(t, &authroutes.SignInRequest{Email: request.Email, Password: request.Password})
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusBadRequest)
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignInRequest()

			data := []struct {
				request func(r authroutes.SignInRequest) authroutes.SignInRequest
				field   string
				name    string
			}{
				{
					field: "email",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Email = ""
						return r
					},
					name: "empty email",
				},
				{
					field: "email",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Email = "   "
						return r
					},
					name: "whitespace email",
				},
				{
					field: "email",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Email = "invalid"
						return r
					},
					name: "invalid email",
				},
				{
					field: "email",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Email = strings.Repeat("t", valueobjects.EmailMaxLength) + "@gmail.com"
						return r
					},
					name: "email exceeds maximum length",
				},
				{
					field: "password",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Password = ""
						return r
					},
					name: "empty password",
				},
				{
					field: "password",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Password = "   "
						return r
					},
					name: "whitespace password",
				},
				{
					field: "password",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Password = strings.Repeat("t", userdomain.PasswordMinLength-1)
						return r
					},
					name: "password too short",
				},
				{
					field: "password",
					request: func(r authroutes.SignInRequest) authroutes.SignInRequest {
						r.Password = strings.Repeat("t", userdomain.PasswordMaxLength+1)
						return r
					},
					name: "password too long",
				},
			}

			for _, data := range data {
				t.Run(data.name, func(t *testing.T) {
					t.Parallel()

					r := data.request(request)
					response := app.SignIn(t, &r)
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
			request := testapp.SignInRequest()
			request.Email = "invalid"

			for i := 0; i < app.Config.RateLimit.SignIn; i++ {
				response := app.SignIn(t, &request)
				defer response.Body.Close()

				app.AssertValidStatus(t, response, http.StatusBadRequest)
			}

			response := app.SignIn(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusTooManyRequests)

		})
	})
}
