package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestResetPassword(t *testing.T) {
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

			forgotPasswordResponse := testApp.ForgotPassword(t, dto.ForgotPasswordRequest{
				Email: validData.Email,
			})
			if forgotPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			resetPasswordToken := testApp.GetResetPasswordToken(t)
			if resetPasswordToken == "" {
				t.Fatal("reste password token should exists after signup")
			}

			resetPasswordResponse := testApp.ResetPassword(t, dto.ResetPasswordRequest{
				Email:       validData.Email,
				Token:       resetPasswordToken,
				NewPassword: strings.Repeat("p", setup.MinPasswordLength),
			})
			if resetPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}
		})
	})

	t.Run("Should Update Password In Database Request Is Valid", func(t *testing.T) {
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

			forgotPasswordResponse := testApp.ForgotPassword(t, dto.ForgotPasswordRequest{
				Email: validData.Email,
			})
			if forgotPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			resetPasswordToken := testApp.GetResetPasswordToken(t)
			if resetPasswordToken == "" {
				t.Fatal("reste password token should exists after signup")
			}

			newPassword := strings.Repeat("p", setup.MinPasswordLength)

			userBeforeUpdate, err := testApp.GetUserByEmail(t.Context(), validData.Email)
			if err != nil {
				t.Fatalf("resetPassword failed: %v", err)
			}

			if userBeforeUpdate == nil {
				t.Fatal("user should exist after signup")
			}

			resetPasswordResponse := testApp.ResetPassword(t, dto.ResetPasswordRequest{
				Email:       validData.Email,
				Token:       resetPasswordToken,
				NewPassword: newPassword,
			})
			if resetPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			userAfterUpdate, err := testApp.GetUserByEmail(t.Context(), validData.Email)
			if err != nil {
				t.Fatalf("resetPassword failed: %v", err)
			}

			if userBeforeUpdate.HashedPassword == userAfterUpdate.HashedPassword {
				t.Fatalf("expected password to be updated, but it remained the same for user %s", userBeforeUpdate.Email)
			}

			if newPassword == userAfterUpdate.HashedPassword {
				t.Fatal("expected the password to be hashed, but it matched the plain text password")
			}
		})
	})

	t.Run("Should Return BadRequest When Request Is Invalid ", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			email := "test@gmail.com"
			newPassword := strings.Repeat("s", setup.MinPasswordLength)

			cases := []struct {
				Email       string
				Token       string
				NewPassword string
				Description string
			}{
				{"", "token", newPassword, "empty email"},
				{"invalid email", "token", newPassword, "invalid email"},
				{email, "", newPassword, "empty token"},
				{email, "token", "", "empty new password"},
				{email, "token", strings.Repeat("s", setup.MinPasswordLength-1), "new password too short"},
				{email, "token", strings.Repeat("s", setup.MaxPasswordLength+1), "new password too long"},
			}

			for _, test := range cases {
				t.Run(test.Description, func(t *testing.T) {
					t.Parallel()

					response := testApp.ResetPassword(t, dto.ResetPasswordRequest{
						Email:       test.Email,
						Token:       test.Token,
						NewPassword: test.NewPassword,
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

			forgotPasswordResponse := testApp.ForgotPassword(t, dto.ForgotPasswordRequest{
				Email: validData.Email,
			})
			if forgotPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			resetPasswordToken := testApp.GetResetPasswordToken(t)
			if resetPasswordToken == "" {
				t.Fatal("reste password token should exists after signup")
			}

			if err := testApp.UnverifyUser(t.Context(), validData.Email); err != nil {
				t.Fatalf("failed to unverify user: %v", err)
			}

			resetPasswordResponse := testApp.ResetPassword(t, dto.ResetPasswordRequest{
				Email:       validData.Email,
				Token:       resetPasswordToken,
				NewPassword: strings.Repeat("p", setup.MinPasswordLength),
			})
			if resetPasswordResponse.StatusCode != http.StatusBadRequest {
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

			verifyResponse := testApp.VerifyAccount(t, dto.VerifyAccountRequest{
				Email: validData.Email,
				Token: verificationToken,
			})
			if verifyResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			forgotPasswordResponse := testApp.ForgotPassword(t, dto.ForgotPasswordRequest{
				Email: validData.Email,
			})
			if forgotPasswordResponse.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			resetPasswordToken := testApp.GetResetPasswordToken(t)
			if resetPasswordToken == "" {
				t.Fatal("reste password token should exists after signup")
			}

			if err := testApp.BanUser(t.Context(), validData.Email); err != nil {
				t.Fatalf("failed to ban user: %v", err)
			}

			resetPasswordResponse := testApp.ResetPassword(t, dto.ResetPasswordRequest{
				Email:       validData.Email,
				Token:       resetPasswordToken,
				NewPassword: strings.Repeat("p", setup.MinPasswordLength),
			})
			if resetPasswordResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
			}

		})
	})
}
