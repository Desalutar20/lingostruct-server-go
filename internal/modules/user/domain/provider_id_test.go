package userdomain

import (
	"strings"
	"testing"
)

func TestNewProviderId(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid provider ID",
			value: "google-user-123",
		},
		{
			name:  "valid provider ID with maximum length",
			value: strings.Repeat("a", ProviderIdMaxLength),
		},
		{
			name:    "empty provider ID",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "provider ID exceeds maximum length",
			value:   strings.Repeat("a", ProviderIdMaxLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProviderId(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewProviderId() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewProviderId() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewProviderId() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}
