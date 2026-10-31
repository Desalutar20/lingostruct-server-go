package tokenutil

import (
	"regexp"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		got, err := Generate(32)
		if err != nil {
			t.Fatalf("Generate() unexpected error: %v", err)
		}

		if len(got.Value()) != 32 {
			t.Errorf("Generate() length = %d, want 32", len(got.Value()))
		}

		if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(got.Value()) {
			t.Errorf("Generate() returned invalid hexadecimal token: %q", got.Value())
		}
	})

	t.Run("odd token length", func(t *testing.T) {
		got, err := Generate(15)
		if err != nil {
			t.Fatalf("Generate() unexpected error: %v", err)
		}

		if len(got.Value()) != 15 {
			t.Errorf("Generate() length = %d, want 15", len(got.Value()))
		}
	})

	t.Run("minimum token length", func(t *testing.T) {
		got, err := Generate(1)
		if err != nil {
			t.Fatalf("Generate() unexpected error: %v", err)
		}

		if len(got.Value()) != 1 {
			t.Errorf("Generate() length = %d, want 1", len(got.Value()))
		}
	})

	t.Run("zero length", func(t *testing.T) {
		_, err := Generate(0)

		if err == nil {
			t.Fatal("Generate() expected an error, got nil")
		}
	})

	t.Run("negative length", func(t *testing.T) {
		_, err := Generate(-1)

		if err == nil {
			t.Fatal("Generate() expected an error, got nil")
		}
	})

	t.Run("unique tokens", func(t *testing.T) {
		first, err := Generate(32)
		if err != nil {
			t.Fatalf("Generate() unexpected error: %v", err)
		}

		second, err := Generate(32)
		if err != nil {
			t.Fatalf("Generate() unexpected error: %v", err)
		}

		if first.Value() == second.Value() {
			t.Error("Generate() returned identical tokens")
		}
	})
}
