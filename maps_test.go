package criterio

import (
	"strings"
	"testing"
)

func TestMapLenMin(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		value   map[string]int
		wantErr bool
	}{
		{"length equals min", 2, map[string]int{"a": 1, "b": 2}, false},
		{"length above min", 2, map[string]int{"a": 1, "b": 2, "c": 3}, false},
		{"length below min", 2, map[string]int{"a": 1}, true},
		{"empty map", 1, map[string]int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapLenMin[string, int](tt.min)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapLenMin(%d) error = %v, wantErr %v", tt.min, err, tt.wantErr)
			}
		})
	}
}

func TestMapLenMax(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		value   map[string]int
		wantErr bool
	}{
		{"length equals max", 3, map[string]int{"a": 1, "b": 2, "c": 3}, false},
		{"length below max", 3, map[string]int{"a": 1, "b": 2}, false},
		{"length above max", 2, map[string]int{"a": 1, "b": 2, "c": 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapLenMax[string, int](tt.max)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapLenMax(%d) error = %v, wantErr %v", tt.max, err, tt.wantErr)
			}
		})
	}
}

func TestMapLenBetween(t *testing.T) {
	tests := []struct {
		name    string
		low     int
		high    int
		value   map[string]int
		wantErr bool
	}{
		{"length equals low", 2, 4, map[string]int{"a": 1, "b": 2}, false},
		{"length equals high", 2, 4, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, false},
		{"length in range", 2, 4, map[string]int{"a": 1, "b": 2, "c": 3}, false},
		{"length below range", 2, 4, map[string]int{"a": 1}, true},
		{"length above range", 2, 4, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapLenBetween[string, int](tt.low, tt.high)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapLenBetween(%d, %d) error = %v, wantErr %v", tt.low, tt.high, err, tt.wantErr)
			}
		})
	}
}

func TestMapNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   map[string]int
		wantErr bool
	}{
		{"non-empty map", map[string]int{"a": 1}, false},
		{"empty map", map[string]int{}, true},
		{"nil map", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapNotEmpty[string, int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	t.Run("all keys pass", func(t *testing.T) {
		validator := MapKeys[string, int](StrNotEmpty)
		m := map[string]int{"hello": 1, "world": 2}
		if err := validator(m); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("key fails validation", func(t *testing.T) {
		validator := MapKeys[string, int](StrMin(3))
		m := map[string]int{"ab": 1, "hello": 2}
		err := validator(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("expected error to mention key, got: %v", err)
		}
	})

	t.Run("empty map passes", func(t *testing.T) {
		validator := MapKeys[string, int](StrNotEmpty)
		if err := validator(map[string]int{}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestMapValues(t *testing.T) {
	t.Run("all values pass", func(t *testing.T) {
		validator := MapValues[string, int](Min(0))
		m := map[string]int{"a": 1, "b": 2}
		if err := validator(m); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("value fails validation", func(t *testing.T) {
		validator := MapValues[string, int](Min(0))
		m := map[string]int{"a": 1, "b": -1}
		err := validator(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "value at key") {
			t.Errorf("expected error to mention value at key, got: %v", err)
		}
	})

	t.Run("empty map passes", func(t *testing.T) {
		validator := MapValues[string, int](Min(0))
		if err := validator(map[string]int{}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("with string values", func(t *testing.T) {
		validator := MapValues[int, string](StrNotEmpty)
		m := map[int]string{1: "hello", 2: "world"}
		if err := validator(m); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
