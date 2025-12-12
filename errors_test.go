package criterio

import (
	"testing"
)

func TestFieldError(t *testing.T) {
	t.Run("formats with field name", func(t *testing.T) {
		err := FieldError{Field: "email", Message: "is invalid"}
		if got := err.Error(); got != "email: is invalid" {
			t.Errorf("got %q, want %q", got, "email: is invalid")
		}
	})

	t.Run("formats without field name", func(t *testing.T) {
		err := FieldError{Message: "is invalid"}
		if got := err.Error(); got != "is invalid" {
			t.Errorf("got %q, want %q", got, "is invalid")
		}
	})
}

func TestFieldErrorsBuilder(t *testing.T) {
	t.Run("empty builder returns nil", func(t *testing.T) {
		var b FieldErrorsBuilder
		if err := b.ToError(); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("append and convert single error", func(t *testing.T) {
		var b FieldErrorsBuilder
		b = b.Append("name", "is required")
		err := b.ToError()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if got := err.Error(); got != "name: is required" {
			t.Errorf("got %q, want %q", got, "name: is required")
		}
	})

	t.Run("append multiple errors", func(t *testing.T) {
		var b FieldErrorsBuilder
		b = b.Append("name", "is required")
		b = b.Append("email", "is invalid")
		err := b.ToError()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "validation failed: name: is required; email: is invalid"
		if got := err.Error(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestFieldErrors(t *testing.T) {
	t.Run("single error format", func(t *testing.T) {
		errs := FieldErrors{{Field: "age", Message: "must be positive"}}
		if got := errs.Error(); got != "age: must be positive" {
			t.Errorf("got %q, want %q", got, "age: must be positive")
		}
	})

	t.Run("multiple errors format", func(t *testing.T) {
		errs := FieldErrors{
			{Field: "name", Message: "is required"},
			{Field: "age", Message: "must be positive"},
		}
		want := "validation failed: name: is required; age: must be positive"
		if got := errs.Error(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("empty errors format", func(t *testing.T) {
		errs := FieldErrors{}
		want := "validation failed with zero error collected"
		if got := errs.Error(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestNewFieldErrors(t *testing.T) {
	errs := NewFieldErrors("field", "message")
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "field" || errs[0].Message != "message" {
		t.Errorf("unexpected error: %+v", errs[0])
	}
}
