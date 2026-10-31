package userdomain

import (
	"strings"
	"testing"
)

func TestLastName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid last name",
			value: "Johnson",
		},
		{
			name:  "valid last name with maximum length",
			value: strings.Repeat("A", LastNameMaxLength),
		},
		{
			name:    "empty last name",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "last name exceeds maximum length",
			value:   strings.Repeat("A", LastNameMaxLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLastName(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewLastName() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewLastName() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewLastName() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}
