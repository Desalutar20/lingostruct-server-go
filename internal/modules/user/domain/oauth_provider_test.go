package userdomain

import "testing"

func TestNewOAuthProvider(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid google provider",
			value: "google",
		},
		{
			name:  "valid github provider",
			value: "github",
		},
		{
			name:    "empty provider",
			value:   "",
			wantErr: true,
		},
		{
			name:    "unsupported provider",
			value:   "facebook",
			wantErr: true,
		},
		{
			name:    "provider with uppercase letters",
			value:   "Google",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOAuthProvider(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NewOAuthProvider() expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("NewOAuthProvider() unexpected error: %v", err)
			}

			if got.Value() != tt.value {
				t.Errorf("NewOAuthProvider() value = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}

func TestOAuthProviderGoogle(t *testing.T) {
	got := OAuthProviderGoogle()

	if got.Value() != "google" {
		t.Errorf("OAuthProviderGoogle() = %q, want %q", got.Value(), "google")
	}
}

func TestOAuthProviderGithub(t *testing.T) {
	got := OAuthProviderGithub()

	if got.Value() != "github" {
		t.Errorf("OAuthProviderGithub() = %q, want %q", got.Value(), "github")
	}
}
