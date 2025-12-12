package criterio

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// StrNotEmpty returns a validator that checks if a string is not empty or whitespace-only.
// Note: This trims whitespace before checking, so "   " is considered empty.
// Use Required[string]() if you only want to reject the zero value "".
func StrNotEmpty() Validator[string] {
	return func(val string) error {
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("cannot be empty")
		}
		return nil
	}
}

// StrMin returns a validator that checks if a string has at least min characters.
// This is a convenience wrapper around MinLen for strings.
func StrMin(min int) Validator[string] {
	return func(val string) error {
		if len(val) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// StrMax returns a validator that checks if a string has at most max characters.
// This is a convenience wrapper around MaxLen for strings.
func StrMax(max int) Validator[string] {
	return func(val string) error {
		if len(val) > max {
			return fmt.Errorf("must be no more than %d characters", max)
		}
		return nil
	}
}

// StrBetween returns a validator that checks if a string length is between low and high (inclusive).
// This is a convenience wrapper around LenBetween for strings.
func StrBetween(low, high int) Validator[string] {
	return func(val string) error {
		if len(val) < low || len(val) > high {
			return fmt.Errorf("must be between %d and %d characters", low, high)
		}
		return nil
	}
}

// StrMatches returns a validator that checks if a string matches the provided regex pattern.
func StrMatches(pattern string) Validator[string] {
	re := regexp.MustCompile(pattern)
	return func(val string) error {
		if !re.MatchString(val) {
			return fmt.Errorf("does not match required pattern")
		}
		return nil
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// StrEmail returns a validator that checks if a string is a valid email address.
// Uses a simple regex pattern for basic validation.
func StrEmail() Validator[string] {
	return func(val string) error {
		if !emailRegex.MatchString(val) {
			return fmt.Errorf("must be a valid email address")
		}
		return nil
	}
}

// StrOneOf returns a validator that checks if a string is one of the allowed values.
// Automatically uses slice iteration for small sets (≤10) and map lookup for larger sets.
func StrOneOf(allowed ...string) Validator[string] {
	// Use slice iteration for small sets (faster due to cache locality)
	if len(allowed) <= 10 {
		return func(val string) error {
			if slices.Contains(allowed, val) {
				return nil
			}
			return fmt.Errorf("must be one of: %s", strings.Join(allowed, ", "))
		}
	}

	// Use map lookup for larger sets (O(1) lookup)
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}
	return func(val string) error {
		if _, ok := allowedSet[val]; !ok {
			return fmt.Errorf("must be one of: %s", strings.Join(allowed, ", "))
		}
		return nil
	}
}
