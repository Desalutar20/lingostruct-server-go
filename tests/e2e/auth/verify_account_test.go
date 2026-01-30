package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestVerifyAccount(t *testing.T) {
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

			verificationResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: validData.Email,
				Token: verificationToken,
			})
			if verificationResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

		})
	})

	t.Run("User Should Be Verified When Request Is Valid", func(t *testing.T) {
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
				Email: validData.Email,
				Token: verificationToken,
			})

			if verificationResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			user, err := testApp.GetUserByEmail(t.Context(), validData.Email)

			if err != nil {
				t.Fatalf("verifyAccount failed: %v", err)
			}

			if user == nil {
				t.Fatal("user should exist after signup")
			}

			if user.IsVerified == false {
				t.Fatalf("expected user.isVerified %t, got %t", true, false)
			}

		})
	})

	t.Run("Should Return BadRequest When Request Is Invalid ", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			email := "test@gmail.com"
			token := "token"

			cases := []struct {
				Email       string
				Token       string
				Description string
			}{
				{"", token, "empty email"},
				{"invalid email", token, "invalid email"},
				{email, "", "empty token"},
			}

			for _, test := range cases {
				t.Run(test.Description, func(t *testing.T) {
					t.Parallel()

					response := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
						Email: test.Email,
						Token: test.Token,
					})

					if response.StatusCode != http.StatusBadRequest {
						t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
					}
				})
			}
		})
	})

	t.Run("Should Return BadRequest When Token Is Deleted", func(t *testing.T) {
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

			testApp.ClearCache(t)

			verificationResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: validData.Email,
				Token: verificationToken,
			})
			if verificationResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}
		})
	})

	t.Run("Should Return BadRequest When Email Is Different", func(t *testing.T) {
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

			err := testApp.BanUser(t.Context(), validData.Email)
			if err != nil {
				t.Fatalf("failed to ban user: %v", err)
			}

			verificationResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: "random@gmail.com",
				Token: verificationToken,
			})
			if verificationResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}
		})
	})
}
