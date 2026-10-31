package userdomain

import "testing"

func TestNewUserRole(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid admin role",
			value: "admin",
		},
		{
			name:  "valid regular role",
			value: "regular",
		},
		{
			name:    "empty role",
			value:   "",
			wantErr: true,
		},
		{
			name:    "unsupported role",
			value:   "moderator",
			wantErr: true,
		},
		{
			name:    "role with uppercase letters",
			value:   "Admin",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserRole(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewUserRole() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewUserRole() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewUserRole() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}

func TestUserRoleAdmin(t *testing.T) {
	got := UserRoleAdmin()

	if got.Value() != "admin" {
		t.Errorf("UserRoleAdmin() = %q, want %q", got.Value(), "admin")
	}
}

func TestUserRoleRegular(t *testing.T) {
	got := UserRoleRegular()

	if got.Value() != "regular" {
		t.Errorf("UserRoleRegular() = %q, want %q", got.Value(), "regular")
	}
}
