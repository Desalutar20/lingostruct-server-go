package valueobjects

import (
	"strings"
	"testing"
)

func TestEmail(t *testing.T) {
	t.Run("NewEmail", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			wantErr bool
		}{
			{
				name:  "valid email",
				value: "john@example.com",
			},
			{
				name:  "valid email with subdomain",
				value: "john@mail.example.com",
			},
			{
				name:  "valid email with plus sign",
				value: "john+test@example.com",
			},
			{
				name:  "valid email with uppercase letters",
				value: "John.Doe@Example.COM",
			},
			{
				name:    "empty email",
				value:   "",
				wantErr: true,
			},
			{
				name:    "email without at sign",
				value:   "johnexample.com",
				wantErr: true,
			},
			{
				name:    "email without domain",
				value:   "john@",
				wantErr: true,
			},
			{
				name:    "email without local part",
				value:   "@example.com",
				wantErr: true,
			},
			{
				name:    "email without top-level domain",
				value:   "john@example",
				wantErr: true,
			},
			{
				name:    "email with consecutive dots in domain",
				value:   "john@example..com",
				wantErr: true,
			},
			{
				name:    "email with spaces",
				value:   "john doe@example.com",
				wantErr: true,
			},
			{
				name:    "email exceeds maximum length",
				value:   strings.Repeat("a", EmailMaxLength) + "@example.com",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewEmail(tt.value)

				if tt.wantErr {
					if err == nil {
						t.Fatal("NewEmail() expected an error, got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("NewEmail() unexpected error: %v", err)
				}

				if got.Value() != tt.value {
					t.Errorf("NewEmail() value = %q, want %q", got.Value(), tt.value)
				}
			})
		}
	})
}
