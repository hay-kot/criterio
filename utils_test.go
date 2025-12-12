package criterio

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("passes with no validators", func(t *testing.T) {
		err := Run("field", "value")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("passes when all validators pass", func(t *testing.T) {
		err := Run("name", "hello",
			StrNotEmpty(),
			StrMin(3),
		)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns single error", func(t *testing.T) {
		err := Run("age", -5, Min(0))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "age: must be at least 0" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("collects multiple errors", func(t *testing.T) {
		err := Run("value", "",
			StrNotEmpty(),
			StrMin(5),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 2 {
			t.Errorf("expected 2 errors, got %d", len(fieldErrs))
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("creates reusable validator", func(t *testing.T) {
		validateEmail := New("email",
			Required[string](),
			StrEmail(),
		)

		if err := validateEmail("test@example.com"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns errors for invalid input", func(t *testing.T) {
		validateEmail := New("email",
			Required[string](),
			StrEmail(),
		)

		err := validateEmail("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "email: must be a valid email address" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("collects all errors", func(t *testing.T) {
		validateName := New("name",
			Required[string](),
			StrMin(3),
			StrMax(10),
		)

		err := validateName("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 2 {
			t.Errorf("expected 2 errors, got %d", len(fieldErrs))
		}
	})
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"empty string fails", "", true},
		{"non-empty string passes", "hello", false},
		{"zero int fails", 0, true},
		{"non-zero int passes", 42, false},
		{"false bool fails", false, true},
		{"true bool passes", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch v := tt.value.(type) {
			case string:
				err = Required[string]()(v)
			case int:
				err = Required[int]()(v)
			case bool:
				err = Required[bool]()(v)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
