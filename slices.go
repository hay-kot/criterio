package criterio

import (
	"errors"
	"fmt"
)

// Static error messages for slice validators.
var (
	errSliceEmpty  = errors.New("must not be empty")
	errSliceUnique = errors.New("must contain unique elements")
)

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

// SliceNotEmpty returns a validator that checks if a slice has at least one element.
func SliceNotEmpty[T any]() Validator[[]T] {
	return func(val []T) error {
		if len(val) == 0 {
			return errSliceEmpty
		}
		return nil
	}
}

// SliceUnique returns a validator that checks if all elements in a slice are unique.
func SliceUnique[T comparable]() Validator[[]T] {
	return func(val []T) error {
		seen := make(map[T]struct{}, len(val))
		for _, v := range val {
			if _, exists := seen[v]; exists {
				return errSliceUnique
			}
			seen[v] = struct{}{}
		}
		return nil
	}
}

// SliceEach returns a validator that applies a validator to each element in a slice.
func SliceEach[T any](validator Validator[T]) Validator[[]T] {
	return func(val []T) error {
		for i, v := range val {
			if err := validator(v); err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}
		}
		return nil
	}
}
