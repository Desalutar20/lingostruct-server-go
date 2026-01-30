package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestSignIn(t *testing.T) {
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

			verificationToken := testApp.GetVerificationToken(t)
			if verificationToken == "" {
				t.Fatal("verification token should exists after signup")
			}

			verifyResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: validData.Email,
				Token: verificationToken,
			})

			if verifyResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			signInResponse := testApp.SignIn(t, dto.SignInRequest{
				Email:    validData.Email,
				Password: validData.Password,
			})

			if signInResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}
		})
	})

	t.Run("Should Return BadRequest When Request Is Invalid ", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			email := "test@gmail.com"
			password := strings.Repeat("s", setup.MinPasswordLength)

			cases := []struct {
				Email       string
				Password    string
				Description string
			}{
				{"", password, "empty email"},
				{"invalid email", password, "invalid email"},
				{email, "", "empty password"},
				{email, strings.Repeat("s", setup.MinPasswordLength-1), "password too short"},
				{email, strings.Repeat("s", setup.MaxPasswordLength+1), "password too long"},
			}

			for _, test := range cases {
				t.Run(test.Description, func(t *testing.T) {
					t.Parallel()

					response := testApp.SignIn(t, dto.SignInRequest{
						Email:    test.Email,
						Password: test.Password,
					})

					if response.StatusCode != http.StatusBadRequest {
						t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
					}
				})
			}
		})
	})

	t.Run("Should Return BadRequest When User Is Not Verified", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			response := testApp.SignUp(t, validData)
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}

			signInResponse := testApp.SignIn(t, dto.SignInRequest{
				Email:    validData.Email,
				Password: validData.Password,
			})

			if signInResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

		})
	})

	t.Run("Should Return BadRequest When User Is Banned", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			response := testApp.SignUp(t, validData)
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}

			verificationToken := testApp.GetVerificationToken(t)
			if verificationToken == "" {
				t.Fatal("verification token should exists after signup")
			}

			verificationResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: "random@gmail.com",
				Token: verificationToken,
			})
			if verificationResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

			testApp.BanUser(t.Context(), validData.Email)

			signInResponse := testApp.SignIn(t, dto.SignInRequest{
				Email:    validData.Email,
				Password: validData.Password,
			})

			if signInResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

		})
	})
}
