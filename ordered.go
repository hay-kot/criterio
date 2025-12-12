package criterio

import (
	"cmp"
	"fmt"
)

// Min returns a validator that checks if a value is at least min.
// Works with any ordered type (integers, floats, strings, etc.).
func Min[T cmp.Ordered](min T) Validator[T] {
	return func(val T) error {
		if val < min {
			return fmt.Errorf("must be at least %v", min)
		}
		return nil
	}
}

// Max returns a validator that checks if a value is at most max.
// Works with any ordered type (integers, floats, strings, etc.).
func Max[T cmp.Ordered](max T) Validator[T] {
	return func(val T) error {
		if val > max {
			return fmt.Errorf("must be at most %v", max)
		}
		return nil
	}
}

// Between returns a validator that checks if a value is between low and high (inclusive).
// Works with any ordered type (integers, floats, strings, etc.).
func Between[T cmp.Ordered](low, high T) Validator[T] {
	return func(val T) error {
		if val < low || val > high {
			return fmt.Errorf("must be between %v and %v", low, high)
		}
		return nil
	}
}
