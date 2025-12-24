package criterio

import (
	"regexp"
	"testing"
)

func TestStrNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"non-empty string", "hello", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"tabs and newlines", "\t\n", true},
		{"string with spaces", "  hello  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrNotEmpty(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrNotEmpty()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrMin(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		value   string
		wantErr bool
	}{
		{"length equals min", 3, "abc", false},
		{"length above min", 3, "abcd", false},
		{"length below min", 3, "ab", true},
		{"empty string", 1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrMin(tt.min)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrMin(%d)(%q) error = %v, wantErr %v", tt.min, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrMax(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		value   string
		wantErr bool
	}{
		{"length equals max", 5, "hello", false},
		{"length below max", 5, "hi", false},
		{"length above max", 5, "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrMax(tt.max)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrMax(%d)(%q) error = %v, wantErr %v", tt.max, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrBetween(t *testing.T) {
	tests := []struct {
		name    string
		low     int
		high    int
		value   string
		wantErr bool
	}{
		{"length equals low", 3, 5, "abc", false},
		{"length equals high", 3, 5, "abcde", false},
		{"length in range", 3, 5, "abcd", false},
		{"length below range", 3, 5, "ab", true},
		{"length above range", 3, 5, "abcdef", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrBetween(tt.low, tt.high)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrBetween(%d, %d)(%q) error = %v, wantErr %v", tt.low, tt.high, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrMatches(t *testing.T) {
	tests := []struct {
		name    string
		pattern *regexp.Regexp
		value   string
		wantErr bool
	}{
		{"matches simple pattern", regexp.MustCompile(`^[a-z]+$`), "hello", false},
		{"does not match", regexp.MustCompile(`^[a-z]+$`), "Hello123", true},
		{"matches digit pattern", regexp.MustCompile(`^\d{3}-\d{4}$`), "123-4567", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrMatches(tt.pattern)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrMatches(%v)(%q) error = %v, wantErr %v", tt.pattern, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"missing local part", "@example.com", true},
		{"spaces", "test @example.com", true},
		{"invalid TLD", "test@example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrEmail(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrEmail()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrOneOf(t *testing.T) {
	t.Run("value in set", func(t *testing.T) {
		validator := StrOneOf("red", "green", "blue")
		if err := validator("green"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("value not in set", func(t *testing.T) {
		validator := StrOneOf("red", "green", "blue")
		if err := validator("yellow"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("large set uses map lookup", func(t *testing.T) {
		options := make([]string, 20)
		for i := range options {
			options[i] = string(rune('a' + i))
		}
		validator := StrOneOf(options...)

		if err := validator("e"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := validator("z"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestStrURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"valid with query", "https://example.com?q=1", false},
		{"missing scheme", "example.com", true},
		{"missing host", "http://", true},
		{"invalid url", "not a url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrURL(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrURL()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrUUID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid lowercase", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid uppercase", "550E8400-E29B-41D4-A716-446655440000", false},
		{"valid mixed case", "550e8400-E29B-41d4-A716-446655440000", false},
		{"missing dashes", "550e8400e29b41d4a716446655440000", true},
		{"too short", "550e8400-e29b-41d4-a716", true},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544000g", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrUUID(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrUUID()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrAlpha(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"lowercase", "hello", false},
		{"uppercase", "HELLO", false},
		{"mixed case", "Hello", false},
		{"with numbers", "hello123", true},
		{"with spaces", "hello world", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrAlpha(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrAlpha()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrAlphanumeric(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"letters only", "hello", false},
		{"numbers only", "123", false},
		{"mixed", "hello123", false},
		{"with spaces", "hello 123", true},
		{"with symbols", "hello@123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrAlphanumeric(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrAlphanumeric()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrNumeric(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"digits only", "123456", false},
		{"with letters", "123abc", true},
		{"with decimal", "123.45", true},
		{"with spaces", "123 456", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrNumeric(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrNumeric()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrContains(t *testing.T) {
	tests := []struct {
		name    string
		substr  string
		value   string
		wantErr bool
	}{
		{"contains substring", "world", "hello world", false},
		{"substring at start", "hello", "hello world", false},
		{"substring at end", "world", "hello world", false},
		{"does not contain", "foo", "hello world", true},
		{"case sensitive", "World", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrContains(tt.substr)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrContains(%q)(%q) error = %v, wantErr %v", tt.substr, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrHasPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		value   string
		wantErr bool
	}{
		{"has prefix", "hello", "hello world", false},
		{"does not have prefix", "world", "hello world", true},
		{"case sensitive", "Hello", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrHasPrefix(tt.prefix)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrHasPrefix(%q)(%q) error = %v, wantErr %v", tt.prefix, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrHasSuffix(t *testing.T) {
	tests := []struct {
		name    string
		suffix  string
		value   string
		wantErr bool
	}{
		{"has suffix", "world", "hello world", false},
		{"does not have suffix", "hello", "hello world", true},
		{"case sensitive", "World", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrHasSuffix(tt.suffix)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrHasSuffix(%q)(%q) error = %v, wantErr %v", tt.suffix, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrNoWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"no whitespace", "helloworld", false},
		{"with space", "hello world", true},
		{"with tab", "hello\tworld", true},
		{"with newline", "hello\nworld", true},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrNoWhitespace(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrNoWhitespace()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrLowercase(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"all lowercase", "hello", false},
		{"with uppercase", "Hello", true},
		{"all uppercase", "HELLO", true},
		{"numbers allowed", "hello123", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrLowercase(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrLowercase()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestStrUppercase(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"all uppercase", "HELLO", false},
		{"with lowercase", "Hello", true},
		{"all lowercase", "hello", true},
		{"numbers allowed", "HELLO123", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StrUppercase(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrUppercase()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
