package criterio

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// StrNotEmpty returns a validator that checks if a string is not empty.
func StrNotEmpty() Validator[string] {
	return func(field, val string) error {
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("cannot be empty")
		}
		return nil
	}
}

// StrMinLen returns a validator that checks if a string has at least min characters.
// This is a convenience wrapper around MinLen for strings.
func StrMinLen(min int) Validator[string] {
	return MinLen[string](min)
}

// StrMaxLen returns a validator that checks if a string has at most max characters.
// This is a convenience wrapper around MaxLen for strings.
func StrMaxLen(max int) Validator[string] {
	return MaxLen[string](max)
}

// StrBetween returns a validator that checks if a string length is between low and high (inclusive).
// This is a convenience wrapper around LenBetween for strings.
func StrBetween(low, high int) Validator[string] {
	return LenBetween[string](low, high)
}

// StrMatches returns a validator that checks if a string matches the provided regex pattern.
func StrMatches(pattern string) Validator[string] {
	re := regexp.MustCompile(pattern)
	return func(field, val string) error {
		if !re.MatchString(val) {
			return fmt.Errorf("does not match required pattern")
		}
		return nil
	}
}

// StrEmail returns a validator that checks if a string is a valid email address.
// Uses a simple regex pattern for basic validation.
func StrEmail() Validator[string] {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return func(field, val string) error {
		if !re.MatchString(val) {
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
		return func(field, val string) error {
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
	return func(field, val string) error {
		if _, ok := allowedSet[val]; !ok {
			return fmt.Errorf("must be one of: %s", strings.Join(allowed, ", "))
		}
		return nil
	}
}
