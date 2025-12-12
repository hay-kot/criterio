package criterio

import "fmt"

// MapLenMin returns a validator that checks if a map has at least min entries.
func MapLenMin[K comparable, V any](min int) Validator[map[K]V] {
	return func(val map[K]V) error {
		if len(val) < min {
			return fmt.Errorf("must have at least %d entries", min)
		}
		return nil
	}
}

// MapLenMax returns a validator that checks if a map has at most max entries.
func MapLenMax[K comparable, V any](max int) Validator[map[K]V] {
	return func(val map[K]V) error {
		if len(val) > max {
			return fmt.Errorf("must have at most %d entries", max)
		}
		return nil
	}
}

// MapLenBetween returns a validator that checks if a map length is between low and high (inclusive).
func MapLenBetween[K comparable, V any](low, high int) Validator[map[K]V] {
	return func(val map[K]V) error {
		length := len(val)
		if length < low || length > high {
			return fmt.Errorf("must have between %d and %d entries", low, high)
		}
		return nil
	}
}

// MapNotEmpty returns a validator that checks if a map has at least one entry.
func MapNotEmpty[K comparable, V any]() Validator[map[K]V] {
	return func(val map[K]V) error {
		if len(val) == 0 {
			return fmt.Errorf("must not be empty")
		}
		return nil
	}
}

// MapKeys returns a validator that applies a validator to each key in a map.
func MapKeys[K comparable, V any](validator Validator[K]) Validator[map[K]V] {
	return func(val map[K]V) error {
		for k := range val {
			if err := validator(k); err != nil {
				return fmt.Errorf("key %v: %w", k, err)
			}
		}
		return nil
	}
}

// MapValues returns a validator that applies a validator to each value in a map.
func MapValues[K comparable, V any](validator Validator[V]) Validator[map[K]V] {
	return func(val map[K]V) error {
		for k, v := range val {
			if err := validator(v); err != nil {
				return fmt.Errorf("value at key %v: %w", k, err)
			}
		}
		return nil
	}
}
