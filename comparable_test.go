package criterio

import (
	"testing"
)

func TestOneOf(t *testing.T) {
	t.Run("string in allowed set", func(t *testing.T) {
		validator := OneOf("a", "b", "c")
		if err := validator("b"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("string not in allowed set", func(t *testing.T) {
		validator := OneOf("a", "b", "c")
		if err := validator("d"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("int in allowed set", func(t *testing.T) {
		validator := OneOf(1, 2, 3)
		if err := validator(2); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("int not in allowed set", func(t *testing.T) {
		validator := OneOf(1, 2, 3)
		if err := validator(4); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("large set uses map lookup", func(t *testing.T) {
		// More than 10 elements triggers map-based lookup
		allowed := make([]int, 20)
		for i := range allowed {
			allowed[i] = i
		}
		validator := OneOf(allowed...)

		if err := validator(15); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := validator(25); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
