package criterio

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// Default regex patterns used by string validators.
// These can be modified to change validation behavior globally.
//
// WARNING: Modifying these affects all validations application-wide.
// Concurrent modification is not thread-safe. If you need to customize
// patterns, do so at program initialization before any validation occurs.
var (
	EmailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	UUIDRegex         = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	AlphaRegex        = regexp.MustCompile(`^[a-zA-Z]+$`)
	AlphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	NumericRegex      = regexp.MustCompile(`^[0-9]+$`)
	WhitespaceRegex   = regexp.MustCompile(`\s`)
)

// StrNotEmpty validates that a string is not empty or whitespace-only.
// Note: This trims whitespace before checking, so "   " is considered empty.
// Use Required[string]() if you only want to reject the zero value "".
func StrNotEmpty(val string) error {
	if strings.TrimSpace(val) == "" {
		return fmt.Errorf("cannot be empty")
	}
	return nil
}

// StrMin returns a validator that checks if a string has at least min characters (runes).
func StrMin(min int) Validator[string] {
	return func(val string) error {
		if utf8.RuneCountInString(val) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// StrMax returns a validator that checks if a string has at most max characters (runes).
func StrMax(max int) Validator[string] {
	return func(val string) error {
		if utf8.RuneCountInString(val) > max {
			return fmt.Errorf("must be no more than %d characters", max)
		}
		return nil
	}
}

// StrBetween returns a validator that checks if a string length is between low and high (inclusive).
// Counts characters (runes), not bytes.
func StrBetween(low, high int) Validator[string] {
	return func(val string) error {
		length := utf8.RuneCountInString(val)
		if length < low || length > high {
			return fmt.Errorf("must be between %d and %d characters", low, high)
		}
		return nil
	}
}

// StrMatches returns a validator that checks if a string matches the provided regex.
func StrMatches(re *regexp.Regexp) Validator[string] {
	return func(val string) error {
		if !re.MatchString(val) {
			return fmt.Errorf("does not match required pattern")
		}
		return nil
	}
}

// StrEmail validates that a string is a valid email address.
// Uses EmailRegex which can be modified to change validation globally.
func StrEmail(val string) error {
	if !EmailRegex.MatchString(val) {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

// StrOneOf returns a validator that checks if a string is one of the allowed values.
// Automatically uses slice iteration for small sets (≤10) and map lookup for larger sets.
// Panics if no allowed values are provided.
func StrOneOf(allowed ...string) Validator[string] {
	if len(allowed) == 0 {
		panic("StrOneOf: at least one allowed value is required")
	}
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
func StrURL(val string) error {
	u, err := url.Parse(val)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("must be a valid URL")
	}
	return nil
}

// StrUUID validates that a string is a valid UUID format.
// Uses UUIDRegex which can be modified to change validation globally.
func StrUUID(val string) error {
	if !UUIDRegex.MatchString(val) {
		return fmt.Errorf("must be a valid UUID")
	}
	return nil
}

// StrAlpha validates that a string contains only letters.
// Uses AlphaRegex which can be modified to change validation globally.
func StrAlpha(val string) error {
	if !AlphaRegex.MatchString(val) {
		return fmt.Errorf("must contain only letters")
	}
	return nil
}

// StrAlphanumeric validates that a string contains only letters and numbers.
// Uses AlphanumericRegex which can be modified to change validation globally.
func StrAlphanumeric(val string) error {
	if !AlphanumericRegex.MatchString(val) {
		return fmt.Errorf("must contain only letters and numbers")
	}
	return nil
}

// StrNumeric validates that a string contains only digits.
// Uses NumericRegex which can be modified to change validation globally.
func StrNumeric(val string) error {
	if !NumericRegex.MatchString(val) {
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

// StrNoWhitespace validates that a string contains no whitespace.
// Uses WhitespaceRegex which can be modified to change validation globally.
func StrNoWhitespace(val string) error {
	if WhitespaceRegex.MatchString(val) {
		return fmt.Errorf("must not contain whitespace")
	}
	return nil
}

// StrLowercase validates that a string is all lowercase.
func StrLowercase(val string) error {
	if val != strings.ToLower(val) {
		return fmt.Errorf("must be lowercase")
	}
	return nil
}

// StrUppercase validates that a string is all uppercase.
func StrUppercase(val string) error {
	if val != strings.ToUpper(val) {
		return fmt.Errorf("must be uppercase")
	}
	return nil
}
