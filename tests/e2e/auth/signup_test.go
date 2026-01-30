package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/tests/e2e/setup"
)

func TestSignUp(t *testing.T) {
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
		})
	})

	t.Run("Account Should Be Created When Request Is Valid", func(t *testing.T) {
		t.Parallel()
		setup.Run(t, func(testApp *setup.TestApp) {
			response := testApp.SignUp(t, validData)
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}

			user, err := testApp.GetUserByEmail(t.Context(), validData.Email)
			if err != nil {
				t.Fatalf("signUp failed: %v", err)
			}

			if user == nil {
				t.Fatal("user should exist after signup")
			}

			if validData.Password == user.HashedPassword {
				t.Fatal("passwords should not match directly")
			}
		})
	})

	t.Run("Should Return BadRequest When Request Is Invalid ", func(t *testing.T) {
		setup.Run(t, func(testApp *setup.TestApp) {
			email := "test@gmail.com"
			password := strings.Repeat("s", setup.MinPasswordLength)
			firstName := strings.Repeat("s", setup.MinFirstNameLength)
			lastName := strings.Repeat("s", setup.MinLastNameLength)

			cases := []struct {
				Email       string
				Password    string
				Firstname   string
				LastName    string
				Description string
			}{
				{"", password, firstName, lastName, "empty email"},
				{"invalid email", password, firstName, lastName, "invalid email"},
				{email, "", firstName, lastName, "empty password"},
				{email, strings.Repeat("s", setup.MinPasswordLength-1), firstName, lastName, "password too short"},
				{email, strings.Repeat("s", setup.MaxPasswordLength+1), firstName, lastName, "password too long"},
				{email, password, "", lastName, "empty firstName"},
				{email, password, strings.Repeat("s", setup.MinFirstNameLength-1), lastName, "firstName too short"},
				{email, password, strings.Repeat("s", setup.MaxFirstNameLength+1), lastName, "firstName too long"},
				{email, password, firstName, "", "empty lastName"},
				{email, password, firstName, strings.Repeat("s", setup.MinLastNameLength-1), "lastName too short"},
				{email, password, firstName, strings.Repeat("s", setup.MaxLastNameLength+1), "lastName too long"},
			}

			for _, test := range cases {
				t.Run(test.Description, func(t *testing.T) {
					t.Parallel()

					response := testApp.SignUp(t, dto.SignUpRequest{
						Email:     test.Email,
						Password:  test.Password,
						FirstName: test.Firstname,
						LastName:  test.LastName,
					})

					if response.StatusCode != http.StatusBadRequest {
						t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.StatusCode)
					}
				})
			}
		})
	})

	t.Run("Should Return BadRequest When User Already Exists", func(t *testing.T) {
		setup.Run(t, func(testApp *setup.TestApp) {
			t.Parallel()
			response := testApp.SignUp(t, validData)
			if response.StatusCode != http.StatusCreated {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}

			secondResponse := testApp.SignUp(t, validData)
			if secondResponse.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
			}
		})
	})

}
