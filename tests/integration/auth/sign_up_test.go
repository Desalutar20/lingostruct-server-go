package auth_test

import (
	"net/http"
	"strings"
	"testing"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/Desalutar20/lingostruct-server-go/tests/integration/testapp"
)

func TestSignUp(t *testing.T) {
	t.Run("Should Return Created When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()

			response := app.SignUp(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusCreated)

			user := app.GetUserByEmail(t, request.Email)
			if user == nil {
				t.Fatal("expected user to be created")
			}

			if user.HashedPassword == &request.Password {
				t.Fatal("expected user password to be hashed")
			}
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()

			data := []struct {
				request func(r authroutes.SignUpRequest) authroutes.SignUpRequest
				field   string
				name    string
			}{
				{
					field: "email",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Email = ""
						return r
					},
					name: "empty email",
				},
				{
					field: "email",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Email = "   "
						return r
					},
					name: "whitespace email",
				},
				{
					field: "email",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Email = "invalid"
						return r
					},
					name: "invalid email",
				},
				{
					field: "email",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Email = strings.Repeat("t", valueobjects.EmailMaxLength) + "@gmail.com"
						return r
					},
					name: "email exceeds maximum length",
				},
				{
					field: "password",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Password = ""
						return r
					},
					name: "empty password",
				},
				{
					field: "password",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Password = "   "
						return r
					},
					name: "whitespace password",
				},
				{
					field: "password",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Password = strings.Repeat("t", userdomain.PasswordMinLength-1)
						return r
					},
					name: "password too short",
				},
				{
					field: "password",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.Password = strings.Repeat("t", userdomain.PasswordMaxLength+1)
						return r
					},
					name: "password too long",
				},
				{
					field: "firstName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.FirstName = ""
						return r
					},
					name: "empty first name",
				},
				{
					field: "firstName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.FirstName = "   "
						return r
					},
					name: "whitespace first name",
				},
				{
					field: "firstName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.FirstName = strings.Repeat("t", userdomain.FirstNameMaxLength+1)
						return r
					},
					name: "first name too long",
				},
				{
					field: "lastName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.LastName = ""
						return r
					},
					name: "empty last name",
				},
				{
					field: "lastName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.LastName = "   "
						return r
					},
					name: "whitespace last name",
				},
				{
					field: "lastName",
					request: func(r authroutes.SignUpRequest) authroutes.SignUpRequest {
						r.LastName = strings.Repeat("t", userdomain.LastNameMaxLength+1)
						return r
					},
					name: "last name too long",
				},
			}

			for _, data := range data {
				t.Run(data.name, func(t *testing.T) {
					t.Parallel()

					r := data.request(request)
					response := app.SignUp(t, &r)
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
			request := testapp.SignUpRequest()
			request.Email = "invalid"

			for i := 0; i < app.Config.RateLimit.SignUp; i++ {
				response := app.SignUp(t, &request)
				defer response.Body.Close()

				app.AssertValidStatus(t, response, http.StatusBadRequest)
			}

			response := app.SignUp(t, &request)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusTooManyRequests)

		})
	})
}
