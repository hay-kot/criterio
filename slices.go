package criterio

import "fmt"

// SliceLenMin returns a validator that checks if a slice has at least min elements.
func SliceLenMin[T any](min int) Validator[[]T] {
	return func(val []T) error {
		if len(val) < min {
			return fmt.Errorf("must have at least %d elements", min)
		}
		return nil
	}
}

// SliceLenMax returns a validator that checks if a slice has at most max elements.
func SliceLenMax[T any](max int) Validator[[]T] {
	return func(val []T) error {
		if len(val) > max {
			return fmt.Errorf("must have at most %d elements", max)
		}
		return nil
	}
}

// SliceLenBetween returns a validator that checks if a slice length is between low and high (inclusive).
func SliceLenBetween[T any](low, high int) Validator[[]T] {
	return func(val []T) error {
		length := len(val)
		if length < low || length > high {
			return fmt.Errorf("must have between %d and %d elements", low, high)
		}
		return nil
	}
}
