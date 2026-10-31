package userdomain

import (
	"strings"
	"testing"
)

func TestNewFirstName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid first name",
			value: "John",
		},
		{
			name:  "valid first name with maximum length",
			value: strings.Repeat("A", FirstNameMaxLength),
		},
		{
			name:    "empty first name",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "first name exceeds maximum length",
			value:   strings.Repeat("A", FirstNameMaxLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFirstName(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewFirstName() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewFirstName() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewFirstName() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}
