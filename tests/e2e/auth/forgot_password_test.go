package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestForgotPassword(t *testing.T) {
	validData := dto.SignUpRequest{
		Email:     "test@gmail.com",
		Password:  strings.Repeat("s", setup.MinPasswordLength),
		FirstName: strings.Repeat("s", setup.MinFirstNameLength),
		LastName:  strings.Repeat("s", setup.MinLastNameLength),
	}

	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			response := testApp.SignUp(t, validData)
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}

			forgotPasswordResponse := testApp.ForgotPassword(t, dto.ForgotPasswordRequest{
				Email: validData.Email,
			})

			if forgotPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}
		})
	})

	t.Run("Should Return BadRequest When Request Is Invalid ", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {

			cases := []struct {
				Email       string
				Description string
			}{
				{"", "empty email"},
				{"invalid email", "invalid email"},
			}

			for _, test := range cases {
				t.Run(test.Description, func(t *testing.T) {
					t.Parallel()

					response := testApp.SignIn(t, dto.ForgotPasswordRequest{
						Email: test.Email,
					})

					if response.StatusCode != http.StatusBadRequest {
						t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
					}
				})
			}
		})
	})

}
