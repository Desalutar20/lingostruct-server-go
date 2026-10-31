package valueobjects

import "testing"

func TestUrl(t *testing.T) {
	t.Run("NewUrl", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			want    string
			wantErr bool
		}{
			{
				name:  "valid HTTP URL",
				value: "http://example.com",
				want:  "http://example.com",
			},
			{
				name:  "valid HTTPS URL",
				value: "https://example.com",
				want:  "https://example.com",
			},
			{
				name:  "URL with path and query",
				value: "https://example.com/users?page=1",
				want:  "https://example.com/users?page=1",
			},
			{
				name:  "URL with port",
				value: "http://localhost:8080",
				want:  "http://localhost:8080",
			},
			{
				name:  "URL with surrounding whitespace",
				value: "  https://example.com  ",
				want:  "https://example.com",
			},
			{
				name:    "empty URL",
				value:   "",
				wantErr: true,
			},
			{
				name:    "whitespace only",
				value:   "   ",
				wantErr: true,
			},
			{
				name:    "invalid URL",
				value:   "://example.com",
				wantErr: true,
			},
			{
				name:    "unsupported scheme",
				value:   "ftp://example.com",
				wantErr: true,
			},
			{
				name:    "missing scheme",
				value:   "example.com",
				wantErr: true,
			},
			{
				name:    "relative URL",
				value:   "/users/123",
				wantErr: true,
			},
			{
				name:    "URL with only HTTP scheme",
				value:   "http:",
				wantErr: true,
			},
			{
				name:    "URL with only HTTPS scheme",
				value:   "https:",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewUrl(tt.value)

				if tt.wantErr {
					if err == nil {
						t.Fatal("NewUrl() expected an error, got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("NewUrl() unexpected error: %v", err)
				}

				if got.Value() != tt.want {
					t.Errorf("Value() = %q, want %q", got.Value(), tt.want)
				}
			})
		}
	})
}
