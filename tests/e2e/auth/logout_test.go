package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestLogout(t *testing.T) {
	validData := dto.SignUpRequest{
		Email:     "test@gmail.com",
		Password:  strings.Repeat("s", setup.MinPasswordLength),
		FirstName: strings.Repeat("s", setup.MinFirstNameLength),
		LastName:  strings.Repeat("s", setup.MinLastNameLength),
	}

	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			cookie := testApp.CreateAndSignIn(t, validData)

			response := testApp.Logout(t, cookie)
			if response.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

		})
	})

	t.Run("Should Return Unauthorized When User Is Not Logged In", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			response := testApp.Logout(t, nil)
			if response.StatusCode != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
			}
		})
	})

}
