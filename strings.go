package criterio

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// StrNotEmpty validates that a string is not empty or whitespace-only.
// Note: This trims whitespace before checking, so "   " is considered empty.
// Use Required[string]() if you only want to reject the zero value "".
var StrNotEmpty Validator[string] = func(val string) error {
	if strings.TrimSpace(val) == "" {
		return fmt.Errorf("cannot be empty")
	}
	return nil
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

// StrEmail validates that a string is a valid email address.
// Uses a simple regex pattern for basic validation.
var StrEmail Validator[string] = func(val string) error {
	if !emailRegex.MatchString(val) {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
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

// StrURL validates that a string is a valid URL.
// Validates that the string can be parsed as a URL with a scheme and host.
var StrURL Validator[string] = func(val string) error {
	u, err := url.Parse(val)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("must be a valid URL")
	}
	return nil
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// StrUUID validates that a string is a valid UUID format.
var StrUUID Validator[string] = func(val string) error {
	if !uuidRegex.MatchString(val) {
		return fmt.Errorf("must be a valid UUID")
	}
	return nil
}

var alphaRegex = regexp.MustCompile(`^[a-zA-Z]+$`)

// StrAlpha validates that a string contains only letters.
var StrAlpha Validator[string] = func(val string) error {
	if !alphaRegex.MatchString(val) {
		return fmt.Errorf("must contain only letters")
	}
	return nil
}

var alphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// StrAlphanumeric validates that a string contains only letters and numbers.
var StrAlphanumeric Validator[string] = func(val string) error {
	if !alphanumericRegex.MatchString(val) {
		return fmt.Errorf("must contain only letters and numbers")
	}
	return nil
}

var numericRegex = regexp.MustCompile(`^[0-9]+$`)

// StrNumeric validates that a string contains only digits.
var StrNumeric Validator[string] = func(val string) error {
	if !numericRegex.MatchString(val) {
		return fmt.Errorf("must contain only digits")
	}
	return nil
}

// StrContains returns a validator that checks if a string contains a substring.
func StrContains(substr string) Validator[string] {
	return func(val string) error {
		if !strings.Contains(val, substr) {
			return fmt.Errorf("must contain %q", substr)
		}
		return nil
	}
}

// StrHasPrefix returns a validator that checks if a string starts with a prefix.
func StrHasPrefix(prefix string) Validator[string] {
	return func(val string) error {
		if !strings.HasPrefix(val, prefix) {
			return fmt.Errorf("must start with %q", prefix)
		}
		return nil
	}
}

// StrHasSuffix returns a validator that checks if a string ends with a suffix.
func StrHasSuffix(suffix string) Validator[string] {
	return func(val string) error {
		if !strings.HasSuffix(val, suffix) {
			return fmt.Errorf("must end with %q", suffix)
		}
		return nil
	}
}

var whitespaceRegex = regexp.MustCompile(`\s`)

// StrNoWhitespace validates that a string contains no whitespace.
var StrNoWhitespace Validator[string] = func(val string) error {
	if whitespaceRegex.MatchString(val) {
		return fmt.Errorf("must not contain whitespace")
	}
	return nil
}

// StrLowercase validates that a string is all lowercase.
var StrLowercase Validator[string] = func(val string) error {
	if val != strings.ToLower(val) {
		return fmt.Errorf("must be lowercase")
	}
	return nil
}

// StrUppercase validates that a string is all uppercase.
var StrUppercase Validator[string] = func(val string) error {
	if val != strings.ToUpper(val) {
		return fmt.Errorf("must be uppercase")
	}
	return nil
}
