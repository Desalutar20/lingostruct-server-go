package valueobjects

import (
	"errors"
	"testing"
)

func TestNonEmptyString(t *testing.T) {
	t.Run("NewNonEmptyString", func(t *testing.T) {
		tests := []struct {
			name      string
			value     string
			minLength int
			maxLength int
			want      string
			wantErr   bool
		}{
			{
				name:  "valid string",
				value: "hello",
				want:  "hello",
			},
			{
				name:  "valid string with surrounding whitespace",
				value: "  hello  ",
				want:  "hello",
			},
			{
				name:    "empty string",
				value:   "",
				wantErr: true,
			},
			{
				name:    "whitespace only",
				value:   "   ",
				wantErr: true,
			},
			{
				name:      "valid string with minimum length",
				value:     "hello",
				minLength: 5,
				want:      "hello",
			},
			{
				name:      "string below minimum length",
				value:     "hell",
				minLength: 5,
				wantErr:   true,
			},
			{
				name:      "valid string with maximum length",
				value:     "hello",
				maxLength: 5,
				want:      "hello",
			},
			{
				name:      "string exceeds maximum length",
				value:     "hello!",
				maxLength: 5,
				wantErr:   true,
			},
			{
				name:      "valid string with minimum and maximum length",
				value:     "hello",
				minLength: 3,
				maxLength: 5,
				want:      "hello",
			},
			{
				name:      "unicode string",
				value:     "Привет",
				minLength: 6,
				maxLength: 6,
				want:      "Привет",
			},
			{
				name:      "string with surrounding whitespace respects trimmed length",
				value:     "  hello  ",
				minLength: 5,
				maxLength: 5,
				want:      "hello",
			},
			{
				name:  "string with tabs and newlines",
				value: "\t hello \n",
				want:  "hello",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewNonEmptyString(
					tt.value,
					"value",
					"Value",
					tt.minLength,
					tt.maxLength,
					nil,
				)

				if tt.wantErr {
					if err == nil {
						t.Fatal("NewNonEmptyString() expected an error, got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("NewNonEmptyString() unexpected error: %v", err)
				}

				if got.Value() != tt.want {
					t.Errorf("Value() = %q, want %q", got.Value(), tt.want)
				}
			})
		}
	})

	t.Run("AdditionalCheck", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			got, err := NewNonEmptyString(
				"hello",
				"value",
				"Value",
				0,
				0,
				func(value string) error {
					if value != "hello" {
						t.Errorf("additionalCheck() value = %q, want %q", value, "hello")
					}
					return nil
				},
			)

			if err != nil {
				t.Fatalf("NewNonEmptyString() unexpected error: %v", err)
			}

			if got.Value() != "hello" {
				t.Errorf("Value() = %q, want %q", got.Value(), "hello")
			}
		})

		t.Run("error", func(t *testing.T) {
			got, err := NewNonEmptyString(
				"hello",
				"value",
				"Value",
				0,
				0,
				func(value string) error {
					return errors.New("")
				},
			)

			if err == nil {
				t.Fatal("NewNonEmptyString() expected an error, got nil")
			}

			if got.Value() != "" {
				t.Errorf("Value() = %q, want empty string", got.Value())
			}
		})
	})

	t.Run("Value", func(t *testing.T) {
		value := NonEmptyString{value: "hello"}

		if got := value.Value(); got != "hello" {
			t.Errorf("Value() = %q, want %q", got, "hello")
		}
	})

	t.Run("String", func(t *testing.T) {
		value := NonEmptyString{value: "hello"}

		if got := value.String(); got != "hello" {
			t.Errorf("String() = %q, want %q", got, "hello")
		}
	})

	t.Run("Compare", func(t *testing.T) {
		tests := []struct {
			name string
			a    string
			b    string
			want int
		}{
			{
				name: "equal strings",
				a:    "hello",
				b:    "hello",
				want: 0,
			},
			{
				name: "first string is less",
				a:    "apple",
				b:    "banana",
				want: -1,
			},
			{
				name: "first string is greater",
				a:    "banana",
				b:    "apple",
				want: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a := NonEmptyString{value: tt.a}
				b := NonEmptyString{value: tt.b}

				got := a.Compare(b)

				if got != tt.want {
					t.Errorf("Compare() = %d, want %d", got, tt.want)
				}
			})
		}
	})
}
