package option

import "testing"

func TestOption(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		got := Some(42)

		if !got.IsSome() {
			t.Error("IsSome() = false, want true")
		}

		if got.IsNone() {
			t.Error("IsNone() = true, want false")
		}

		if value := got.Unwrap(); value != 42 {
			t.Errorf("Unwrap() = %d, want 42", value)
		}

		if value := got.UnwrapOr(10); value != 42 {
			t.Errorf("UnwrapOr() = %d, want 42", value)
		}
	})

	t.Run("None", func(t *testing.T) {
		got := None[int]()

		if got.IsSome() {
			t.Error("IsSome() = true, want false")
		}

		if !got.IsNone() {
			t.Error("IsNone() = false, want true")
		}

		if value := got.Unwrap(); value != 0 {
			t.Errorf("Unwrap() = %d, want 0", value)
		}

		if value := got.UnwrapOr(10); value != 10 {
			t.Errorf("UnwrapOr() = %d, want 10", value)
		}
	})

	t.Run("Some with zero value", func(t *testing.T) {
		got := Some(0)

		if !got.IsSome() {
			t.Error("IsSome() = false, want true")
		}

		if got.IsNone() {
			t.Error("IsNone() = true, want false")
		}

		if value := got.UnwrapOr(10); value != 0 {
			t.Errorf("UnwrapOr() = %d, want 0", value)
		}
	})

	t.Run("Some with string", func(t *testing.T) {
		got := Some("hello")

		if value := got.Unwrap(); value != "hello" {
			t.Errorf("Unwrap() = %q, want %q", value, "hello")
		}
	})

	t.Run("None with string", func(t *testing.T) {
		got := None[string]()

		if value := got.UnwrapOr("default"); value != "default" {
			t.Errorf("UnwrapOr() = %q, want %q", value, "default")
		}
	})
}
