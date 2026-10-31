package passwordutil

import (
	"testing"

	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

func TestPasswordUtil(t *testing.T) {
	password, err := userdomain.NewPassword("StrongPassword123")
	if err != nil {
		t.Fatalf("NewPassword() unexpected error: %v", err)
	}

	hashedPassword, err := HashPassword(&password)
	if err != nil {
		t.Fatalf("HashPassword() unexpected error: %v", err)
	}

	t.Run("HashPassword", func(t *testing.T) {
		if hashedPassword.Value() == password.Value() {
			t.Error("HashPassword() returned the original password")
		}

		if len(hashedPassword.Value()) == 0 {
			t.Error("HashPassword() returned an empty hash")
		}

		if len(hashedPassword.Value()) < 40 {
			t.Errorf("HashPassword() hash length = %d, want at least 40", len(hashedPassword.Value()))
		}
	})

	t.Run("VerifyPassword", func(t *testing.T) {
		t.Run("correct password", func(t *testing.T) {
			if !VerifyPassword(&password, &hashedPassword) {
				t.Error("VerifyPassword() = false, want true")
			}
		})

		t.Run("incorrect password", func(t *testing.T) {
			wrongPassword, err := userdomain.NewPassword("WrongPassword123")
			if err != nil {
				t.Fatalf("NewPassword() unexpected error: %v", err)
			}

			if VerifyPassword(&wrongPassword, &hashedPassword) {
				t.Error("VerifyPassword() = true, want false")
			}
		})

		t.Run("different hash for same password", func(t *testing.T) {
			secondHash, err := HashPassword(&password)
			if err != nil {
				t.Fatalf("HashPassword() unexpected error: %v", err)
			}

			if secondHash.Value() == hashedPassword.Value() {
				t.Error("HashPassword() returned the same hash for the same password")
			}

			if !VerifyPassword(&password, &secondHash) {
				t.Error("VerifyPassword() = false, want true for second hash")
			}
		})
	})
}
