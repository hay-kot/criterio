package criterio

import (
	"fmt"
	"slices"
)

// OneOf returns a validator that checks if a value is one of the allowed values.
// Works with any comparable type. Automatically uses slice iteration for small sets
// (≤10) and map lookup for larger sets for optimal performance.
// Panics if no allowed values are provided.
func OneOf[T comparable](allowed ...T) Validator[T] {
	if len(allowed) == 0 {
		panic("OneOf: at least one allowed value is required")
	}
	// Use slice iteration for small sets (faster due to cache locality)
	if len(allowed) <= 10 {
		return func(val T) error {
			if slices.Contains(allowed, val) {
				return nil
			}
			return fmt.Errorf("must be one of the allowed values")
		}
	}

	// Use map lookup for larger sets (O(1) lookup)
	allowedSet := make(map[T]struct{}, len(allowed))
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}
	return func(val T) error {
		if _, ok := allowedSet[val]; !ok {
			return fmt.Errorf("must be one of the allowed values")
		}
		return nil
	}
}
