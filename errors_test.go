package criterio

import (
	"errors"
	"testing"
)

func TestFieldError(t *testing.T) {
	t.Run("formats with field name", func(t *testing.T) {
		err := FieldError{Field: "email", Err: errors.New("is invalid")}
		if got := err.Error(); got != "email: is invalid" {
			t.Errorf("got %q, want %q", got, "email: is invalid")
		}
	})

	t.Run("formats without field name", func(t *testing.T) {
		err := FieldError{Err: errors.New("is invalid")}
		if got := err.Error(); got != "is invalid" {
			t.Errorf("got %q, want %q", got, "is invalid")
		}
	})

	t.Run("unwrap returns underlying error", func(t *testing.T) {
		underlying := errors.New("is invalid")
		err := FieldError{Field: "email", Err: underlying}
		if !errors.Is(err.Unwrap(), underlying) {
			t.Errorf("Unwrap() did not return underlying error")
		}
	})

	t.Run("errors.Is works with wrapped error", func(t *testing.T) {
		underlying := errors.New("is invalid")
		err := FieldError{Field: "email", Err: underlying}
		if !errors.Is(err, underlying) {
			t.Error("errors.Is should match underlying error")
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
		b = b.Append("name", errors.New("is required"))
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
		b = b.Append("name", errors.New("is required"))
		b = b.Append("email", errors.New("is invalid"))
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
		errs := FieldErrors{{Field: "age", Err: errors.New("must be positive")}}
		if got := errs.Error(); got != "age: must be positive" {
			t.Errorf("got %q, want %q", got, "age: must be positive")
		}
	})

	t.Run("multiple errors format", func(t *testing.T) {
		errs := FieldErrors{
			{Field: "name", Err: errors.New("is required")},
			{Field: "age", Err: errors.New("must be positive")},
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
	underlying := errors.New("message")
	errs := NewFieldErrors("field", underlying)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "field" || !errors.Is(errs[0].Err, underlying) {
		t.Errorf("unexpected error: %+v", errs[0])
	}
}
