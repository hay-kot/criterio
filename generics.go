package criterio

import (
	"cmp"
	"fmt"
	"slices"
)

// Lengther is a constraint for types that have a length.
type Lengther interface {
	~string | ~[]byte
}

// MinLen returns a validator that checks if a value has at least min length.
// Works with strings and byte slices.
func MinLen[T Lengther](min int) Validator[T] {
	return func(field string, val T) error {
		if len(val) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// MaxLen returns a validator that checks if a value has at most max length.
// Works with strings and byte slices.
func MaxLen[T Lengther](max int) Validator[T] {
	return func(field string, val T) error {
		if len(val) > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

// LenBetween returns a validator that checks if a value's length is between low and high (inclusive).
// Works with strings and byte slices.
func LenBetween[T Lengther](low, high int) Validator[T] {
	return func(field string, val T) error {
		length := len(val)
		if length < low || length > high {
			return fmt.Errorf("must be between %d and %d characters", low, high)
		}
		return nil
	}
}

// SliceMinLen returns a validator that checks if a slice has at least min elements.
func SliceMinLen[T any](min int) Validator[[]T] {
	return func(field string, val []T) error {
		if len(val) < min {
			return fmt.Errorf("must have at least %d elements", min)
		}
		return nil
	}
}

// SliceMaxLen returns a validator that checks if a slice has at most max elements.
func SliceMaxLen[T any](max int) Validator[[]T] {
	return func(field string, val []T) error {
		if len(val) > max {
			return fmt.Errorf("must have at most %d elements", max)
		}
		return nil
	}
}

// SliceLenBetween returns a validator that checks if a slice length is between low and high (inclusive).
func SliceLenBetween[T any](low, high int) Validator[[]T] {
	return func(field string, val []T) error {
		length := len(val)
		if length < low || length > high {
			return fmt.Errorf("must have between %d and %d elements", low, high)
		}
		return nil
	}
}

// Min returns a validator that checks if a value is at least min.
// Works with any ordered type (integers, floats, strings, etc.).
func Min[T cmp.Ordered](min T) Validator[T] {
	return func(field string, val T) error {
		if val < min {
			return fmt.Errorf("must be at least %v", min)
		}
		return nil
	}
}

// Max returns a validator that checks if a value is at most max.
// Works with any ordered type (integers, floats, strings, etc.).
func Max[T cmp.Ordered](max T) Validator[T] {
	return func(field string, val T) error {
		if val > max {
			return fmt.Errorf("must be at most %v", max)
		}
		return nil
	}
}

// Between returns a validator that checks if a value is between low and high (inclusive).
// Works with any ordered type (integers, floats, strings, etc.).
func Between[T cmp.Ordered](low, high T) Validator[T] {
	return func(field string, val T) error {
		if val < low || val > high {
			return fmt.Errorf("must be between %v and %v", low, high)
		}
		return nil
	}
}

// OneOf returns a validator that checks if a value is one of the allowed values.
// Works with any comparable type. Automatically uses slice iteration for small sets
// (≤10) and map lookup for larger sets for optimal performance.
func OneOf[T comparable](allowed ...T) Validator[T] {
	// Use slice iteration for small sets (faster due to cache locality)
	if len(allowed) <= 10 {
		return func(field string, val T) error {
			if slices.Contains(allowed, val) {
				return nil
			}
			return fmt.Errorf("must be one of the allowed values: %s")
		}
	}

	// Use map lookup for larger sets (O(1) lookup)
	allowedSet := make(map[T]struct{}, len(allowed))
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}
	return func(field string, val T) error {
		if _, ok := allowedSet[val]; !ok {
			return fmt.Errorf("must be one of the allowed values")
		}
		return nil
	}
}
