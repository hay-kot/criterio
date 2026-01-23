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
			StrNotEmpty,
			StrMin(3),
		)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns first error only", func(t *testing.T) {
		err := Run("age", -5, Min(0))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "age: must be at least 0" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("stops on first failure", func(t *testing.T) {
		err := Run("value", "",
			StrNotEmpty,
			StrMin(5),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 1 {
			t.Errorf("expected 1 error (fail-fast), got %d", len(fieldErrs))
		}
	})
}

func TestRunAll(t *testing.T) {
	t.Run("passes with no validators", func(t *testing.T) {
		err := RunAll("field", "value")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("passes when all validators pass", func(t *testing.T) {
		err := RunAll("name", "hello",
			StrNotEmpty,
			StrMin(3),
		)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("collects multiple errors", func(t *testing.T) {
		err := RunAll("value", "",
			StrNotEmpty,
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
			Required[string],
			StrEmail,
		)

		if err := validateEmail("test@example.com"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns first error for invalid input", func(t *testing.T) {
		validateEmail := New("email",
			Required[string],
			StrEmail,
		)

		err := validateEmail("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "email: must be a valid email address" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("stops on first failure", func(t *testing.T) {
		validateName := New("name",
			Required[string],
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
		if len(fieldErrs) != 1 {
			t.Errorf("expected 1 error (fail-fast), got %d", len(fieldErrs))
		}
	})
}

func TestNewAll(t *testing.T) {
	t.Run("creates reusable validator", func(t *testing.T) {
		validateEmail := NewAll("email",
			Required[string],
			StrEmail,
		)

		if err := validateEmail("test@example.com"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("collects all errors", func(t *testing.T) {
		validateName := NewAll("name",
			Required[string],
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
				err = Required[string](v)
			case int:
				err = Required[int](v)
			case bool:
				err = Required[bool](v)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	t.Run("returns nil when all validations pass", func(t *testing.T) {
		err := ValidateStruct(
			Run("name", "John", StrNotEmpty),
			Run("age", 25, Min(0)),
		)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("collects errors from multiple fields", func(t *testing.T) {
		err := ValidateStruct(
			Run("name", "", StrNotEmpty),
			Run("age", -5, Min(0)),
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

	t.Run("handles nil errors", func(t *testing.T) {
		err := ValidateStruct(
			nil,
			Run("name", "John", StrNotEmpty),
			nil,
		)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("works with RunAll", func(t *testing.T) {
		err := ValidateStruct(
			RunAll("name", "", StrNotEmpty, StrMin(3)),
			Run("age", -5, Min(0)),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 3 {
			t.Errorf("expected 3 errors, got %d", len(fieldErrs))
		}
	})

	t.Run("handles non-FieldErrors error", func(t *testing.T) {
		err := ValidateStruct(
			errors.New("generic error"),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if fieldErrs[0].Field != "" {
			t.Errorf("expected empty field, got %s", fieldErrs[0].Field)
		}
		if fieldErrs[0].Message != "generic error" {
			t.Errorf("expected 'generic error', got %s", fieldErrs[0].Message)
		}
	})
}

func TestNest(t *testing.T) {
	t.Run("returns nil for nil error", func(t *testing.T) {
		err := Nest("address", nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("prefixes single field error", func(t *testing.T) {
		err := Nest("address", Run("street", "", StrNotEmpty))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "address.street: cannot be empty" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("prefixes multiple field errors", func(t *testing.T) {
		innerErr := ValidateStruct(
			Run("street", "", StrNotEmpty),
			Run("city", "", StrNotEmpty),
		)
		err := Nest("address", innerErr)
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
		if fieldErrs[0].Field != "address.street" {
			t.Errorf("expected address.street, got %s", fieldErrs[0].Field)
		}
		if fieldErrs[1].Field != "address.city" {
			t.Errorf("expected address.city, got %s", fieldErrs[1].Field)
		}
	})

	t.Run("handles double nesting", func(t *testing.T) {
		innerErr := Run("zip", "", StrNotEmpty)
		nestedOnce := Nest("address", innerErr)
		nestedTwice := Nest("user", nestedOnce)
		if nestedTwice == nil {
			t.Fatal("expected error, got nil")
		}
		if nestedTwice.Error() != "user.address.zip: cannot be empty" {
			t.Errorf("unexpected error: %v", nestedTwice)
		}
	})

	t.Run("handles triple nesting", func(t *testing.T) {
		// Simulate: company.office.address.street
		innerErr := Run("street", "", StrNotEmpty)
		level1 := Nest("address", innerErr)
		level2 := Nest("office", level1)
		level3 := Nest("company", level2)
		if level3 == nil {
			t.Fatal("expected error, got nil")
		}
		if level3.Error() != "company.office.address.street: cannot be empty" {
			t.Errorf("unexpected error: %v", level3)
		}
	})

	t.Run("handles deeply nested struct pattern", func(t *testing.T) {
		// Simulates nested Validate() calls with multiple fields
		addressErr := ValidateStruct(
			Run("street", "", StrNotEmpty),
			Run("zip", "abc", StrNumeric),
		)
		officeErr := ValidateStruct(
			Run("name", "", StrNotEmpty),
			Nest("address", addressErr),
		)
		companyErr := ValidateStruct(
			Nest("headquarters", officeErr),
		)

		var fieldErrs FieldErrors
		if !errors.As(companyErr, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 3 {
			t.Errorf("expected 3 errors, got %d", len(fieldErrs))
		}

		expectedFields := []string{
			"headquarters.name",
			"headquarters.address.street",
			"headquarters.address.zip",
		}
		for i, expected := range expectedFields {
			if fieldErrs[i].Field != expected {
				t.Errorf("expected %s, got %s", expected, fieldErrs[i].Field)
			}
		}
	})

	t.Run("handles empty field in original error", func(t *testing.T) {
		err := Nest("config", FieldErrors{{Field: "", Message: "invalid format"}})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if fieldErrs[0].Field != "config" {
			t.Errorf("expected config, got %s", fieldErrs[0].Field)
		}
	})

	t.Run("works in ValidateStruct", func(t *testing.T) {
		addressErr := ValidateStruct(
			Run("street", "", StrNotEmpty),
		)
		err := ValidateStruct(
			Run("name", "John", StrNotEmpty),
			Nest("address", addressErr),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "address.street: cannot be empty" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("handles non-FieldErrors error", func(t *testing.T) {
		err := Nest("config", errors.New("invalid format"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "config: invalid format" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestWhen(t *testing.T) {
	t.Run("runs validators when condition is true", func(t *testing.T) {
		validator := When(true, StrNotEmpty)
		if err := validator(""); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("skips validators when condition is false", func(t *testing.T) {
		validator := When(false, StrNotEmpty)
		if err := validator(""); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("chains multiple validators", func(t *testing.T) {
		validator := When(true, StrNotEmpty, StrMin(5))
		if err := validator("hi"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("passes when condition true and value valid", func(t *testing.T) {
		validator := When(true, StrNotEmpty, StrMin(3))
		if err := validator("hello"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("panics with empty validators", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got nil")
			}
		}()
		When[string](true)
	})
}

func TestSkipIf(t *testing.T) {
	t.Run("skips validators when condition is true", func(t *testing.T) {
		validator := SkipIf(true, StrNotEmpty)
		if err := validator(""); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("runs validators when condition is false", func(t *testing.T) {
		validator := SkipIf(false, StrNotEmpty)
		if err := validator(""); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestOr(t *testing.T) {
	t.Run("passes if first validator passes", func(t *testing.T) {
		validator := Or(StrEmail, StrNumeric)
		if err := validator("test@example.com"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("passes if second validator passes", func(t *testing.T) {
		validator := Or(StrEmail, StrNumeric)
		if err := validator("12345"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("fails if all validators fail", func(t *testing.T) {
		validator := Or(StrEmail, StrNumeric)
		if err := validator("not-email-or-numeric"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns last error on failure", func(t *testing.T) {
		validator := Or(StrEmail, StrNumeric)
		err := validator("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "must contain only digits" {
			t.Errorf("expected last validator error, got %v", err)
		}
	})

	t.Run("panics with empty validators", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got nil")
			}
		}()
		Or[string]()
	})
}

func TestNot(t *testing.T) {
	t.Run("passes when wrapped validator fails", func(t *testing.T) {
		validator := Not(StrNumeric, "must not be numeric")
		if err := validator("hello"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("fails when wrapped validator passes", func(t *testing.T) {
		validator := Not(StrNumeric, "must not be numeric")
		err := validator("12345")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "must not be numeric" {
			t.Errorf("expected custom message, got %v", err)
		}
	})
}

func TestValidateSlice(t *testing.T) {
	type Address struct {
		Street string
		Zip    string
	}

	validateAddress := func(a Address) error {
		return ValidateStruct(
			Run("street", a.Street, Required[string]),
			Run("zip", a.Zip, Required[string], StrNumeric),
		)
	}

	t.Run("returns nil for empty slice", func(t *testing.T) {
		err := ValidateSlice("addresses", []Address{}, validateAddress)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns nil when all elements valid", func(t *testing.T) {
		addresses := []Address{
			{Street: "123 Main St", Zip: "12345"},
			{Street: "456 Oak Ave", Zip: "67890"},
		}
		err := ValidateSlice("addresses", addresses, validateAddress)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("returns indexed errors for invalid elements", func(t *testing.T) {
		addresses := []Address{
			{Street: "123 Main St", Zip: "12345"},
			{Street: "", Zip: "abc"},
		}
		err := ValidateSlice("addresses", addresses, validateAddress)
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
		if fieldErrs[0].Field != "addresses[1].street" {
			t.Errorf("expected addresses[1].street, got %s", fieldErrs[0].Field)
		}
		if fieldErrs[1].Field != "addresses[1].zip" {
			t.Errorf("expected addresses[1].zip, got %s", fieldErrs[1].Field)
		}
	})

	t.Run("collects errors from multiple elements", func(t *testing.T) {
		addresses := []Address{
			{Street: "", Zip: "12345"},
			{Street: "456 Oak Ave", Zip: ""},
		}
		err := ValidateSlice("addresses", addresses, validateAddress)
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
		if fieldErrs[0].Field != "addresses[0].street" {
			t.Errorf("expected addresses[0].street, got %s", fieldErrs[0].Field)
		}
		if fieldErrs[1].Field != "addresses[1].zip" {
			t.Errorf("expected addresses[1].zip, got %s", fieldErrs[1].Field)
		}
	})

	t.Run("works with simple validators", func(t *testing.T) {
		emails := []string{"valid@example.com", "invalid", "also@valid.com"}
		err := ValidateSlice("emails", emails, StrEmail)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 1 {
			t.Errorf("expected 1 error, got %d", len(fieldErrs))
		}
		if fieldErrs[0].Field != "emails[1]" {
			t.Errorf("expected emails[1], got %s", fieldErrs[0].Field)
		}
	})

	t.Run("works in ValidateStruct", func(t *testing.T) {
		type User struct {
			Name      string
			Addresses []Address
		}

		user := User{
			Name: "John",
			Addresses: []Address{
				{Street: "", Zip: "12345"},
			},
		}

		err := ValidateStruct(
			Run("name", user.Name, Required[string]),
			ValidateSlice("addresses", user.Addresses, validateAddress),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var fieldErrs FieldErrors
		if !errors.As(err, &fieldErrs) {
			t.Fatal("expected FieldErrors type")
		}
		if len(fieldErrs) != 1 {
			t.Errorf("expected 1 error, got %d", len(fieldErrs))
		}
		if fieldErrs[0].Field != "addresses[0].street" {
			t.Errorf("expected addresses[0].street, got %s", fieldErrs[0].Field)
		}
	})
}
