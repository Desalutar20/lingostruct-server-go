package userdomain

import (
	"strings"
	"testing"
)

func TestNewHashedPassword(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid hashed password",
			value: strings.Repeat("a", HashedPasswordMinLength+10),
		},
		{
			name:  "valid hashed password with minimum length",
			value: strings.Repeat("a", HashedPasswordMinLength),
		},
		{
			name:  "valid hashed password with maximum length",
			value: strings.Repeat("a", HashedPasswordMaxLength),
		},
		{
			name:    "empty hashed password",
			value:   "",
			wantErr: true,
		},
		{
			name:    "hashed password below minimum length",
			value:   strings.Repeat("a", HashedPasswordMinLength-1),
			wantErr: true,
		},
		{
			name:    "hashed password exceeds maximum length",
			value:   strings.Repeat("a", HashedPasswordMaxLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHashedPassword(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewHashedPassword() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewHashedPassword() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewHashedPassword() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}
