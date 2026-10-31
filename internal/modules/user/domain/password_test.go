package userdomain

import (
	"strings"
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid password",
			value: "password123",
		},
		{
			name:  "valid password with minimum length",
			value: strings.Repeat("a", PasswordMinLength),
		},
		{
			name:  "valid password with maximum length",
			value: strings.Repeat("a", PasswordMaxLength),
		},
		{
			name:    "empty password",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			value:   "        ",
			wantErr: true,
		},
		{
			name:    "password below minimum length",
			value:   strings.Repeat("a", PasswordMinLength-1),
			wantErr: true,
		},
		{
			name:    "password exceeds maximum length",
			value:   strings.Repeat("a", PasswordMaxLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPassword(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewPassword() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewPassword() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewPassword() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}
