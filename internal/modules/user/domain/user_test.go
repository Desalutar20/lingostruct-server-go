package userdomain

import (
	"strings"
	"testing"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/google/uuid"
)

func newTestUser(t *testing.T) User {
	t.Helper()

	email, err := valueobjects.NewEmail("john@example.com")
	if err != nil {
		t.Fatalf("NewEmail() unexpected error: %v", err)
	}

	password, err := NewHashedPassword(strings.Repeat("a", HashedPasswordMinLength))
	if err != nil {
		t.Fatalf("NewHashedPassword() unexpected error: %v", err)
	}

	return NewUser(email, &password, nil, nil)
}

func TestUser(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		user := newTestUser(t)

		if user.FirstName() != nil {
			t.Error("FirstName() should be nil")
		}

		if user.LastName() != nil {
			t.Error("LastName() should be nil")
		}

		if user.HashedPassword() == nil {
			t.Error("HashedPassword() should not be nil")
		}

		if user.Role() != UserRoleRegular() {
			t.Errorf("Role() = %q, want %q", user.Role().Value(), "regular")
		}

		if user.IsBanned() {
			t.Error("IsBanned() = true, want false")
		}

		if user.IsVerified() {
			t.Error("IsVerified() = true, want false")
		}

		if user.AvatarUrl() != nil {
			t.Error("AvatarUrl() should be nil")
		}

		if user.GoogleId() != nil {
			t.Error("GoogleId() should be nil")
		}

		if user.GithubId() != nil {
			t.Error("GithubId() should be nil")
		}
	})

	t.Run("Verify", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			user := newTestUser(t)

			err := user.Verify()
			if err != nil {
				t.Fatalf("Verify() unexpected error: %v", err)
			}

			if !user.IsVerified() {
				t.Error("IsVerified() = false, want true")
			}
		})

		t.Run("banned user", func(t *testing.T) {
			user := RestoreUser(
				uuid.Nil, time.Time{}, time.Time{},
				"john@example.com",
				nil, nil, nil,
				"regular",
				true, false,
				nil, nil, nil,
			)

			err := user.Verify()
			if err == nil {
				t.Fatal("Verify() expected an error, got nil")
			}

			if user.IsVerified() {
				t.Error("IsVerified() = true, want false")
			}
		})
	})

	t.Run("UpdatePassword", func(t *testing.T) {
		user := newTestUser(t)

		password, err := NewHashedPassword(strings.Repeat("A", HashedPasswordMinLength))
		if err != nil {
			t.Fatalf("NewHashedPassword() unexpected error: %v", err)
		}

		user.UpdatePassword(password)

		got := user.HashedPassword()
		if got == nil {
			t.Fatal("HashedPassword() should not be nil")
		}

		if got.Value() != password.Value() {
			t.Errorf("HashedPassword() = %q, want %q", got.Value(), password.Value())
		}
	})

	t.Run("LinkOAuthProvider", func(t *testing.T) {
		t.Run("google", func(t *testing.T) {
			user := newTestUser(t)

			providerID, err := NewProviderId("google-user-123")
			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			if err := user.LinkOAuthProvider(OAuthProviderGoogle(), providerID); err != nil {
				t.Fatalf("LinkOAuthProvider() unexpected error: %v", err)
			}

			if got := user.GoogleId(); got == nil || got.Value() != providerID.Value() {
				t.Errorf("GoogleId() = %v, want %q", got, providerID.Value())
			}
		})

		t.Run("github", func(t *testing.T) {
			user := newTestUser(t)

			providerID, err := NewProviderId("github-user-123")
			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			if err := user.LinkOAuthProvider(OAuthProviderGithub(), providerID); err != nil {
				t.Fatalf("LinkOAuthProvider() unexpected error: %v", err)
			}

			if got := user.GithubId(); got == nil || got.Value() != providerID.Value() {
				t.Errorf("GithubId() = %v, want %q", got, providerID.Value())
			}
		})

		t.Run("already linked provider", func(t *testing.T) {
			user := newTestUser(t)

			firstID, err := NewProviderId("google-user-123")
			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			secondID, err := NewProviderId("google-user-456")
			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			if err := user.LinkOAuthProvider(OAuthProviderGoogle(), firstID); err != nil {
				t.Fatalf("LinkOAuthProvider() unexpected error: %v", err)
			}

			if err := user.LinkOAuthProvider(OAuthProviderGoogle(), secondID); err != nil {
				t.Fatalf("LinkOAuthProvider() unexpected error: %v", err)
			}

			if got := user.GoogleId().Value(); got != firstID.Value() {
				t.Errorf("GoogleId() = %q, want %q", got, firstID.Value())
			}
		})

		t.Run("unsupported provider", func(t *testing.T) {
			user := newTestUser(t)

			provider := OAuthProvider{
				NonEmptyString: valueobjects.NonEmptyString{},
			}

			providerID, err := NewProviderId("provider-123")
			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			if err := user.LinkOAuthProvider(provider, providerID); err == nil {
				t.Fatal("LinkOAuthProvider() expected an error, got nil")
			}
		})
	})
}
