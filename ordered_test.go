package criterio

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		value   int
		wantErr bool
	}{
		{"value equals min", 5, 5, false},
		{"value above min", 5, 10, false},
		{"value below min", 5, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Min(tt.min)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Min(%d)(%d) error = %v, wantErr %v", tt.min, tt.value, err, tt.wantErr)
			}
		})
	}

	t.Run("works with floats", func(t *testing.T) {
		if err := Min(1.5)(2.0); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := Min(1.5)(1.0); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		if err := Min("b")("c"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := Min("b")("a"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMax(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		value   int
		wantErr bool
	}{
		{"value equals max", 10, 10, false},
		{"value below max", 10, 5, false},
		{"value above max", 10, 15, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Max(tt.max)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Max(%d)(%d) error = %v, wantErr %v", tt.max, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		name    string
		low     int
		high    int
		value   int
		wantErr bool
	}{
		{"value equals low", 5, 10, 5, false},
		{"value equals high", 5, 10, 10, false},
		{"value in range", 5, 10, 7, false},
		{"value below range", 5, 10, 3, true},
		{"value above range", 5, 10, 15, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Between(tt.low, tt.high)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Between(%d, %d)(%d) error = %v, wantErr %v", tt.low, tt.high, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestPositive(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"positive value", 5, false},
		{"zero", 0, true},
		{"negative value", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Positive[int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Positive()(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}

	t.Run("works with floats", func(t *testing.T) {
		if err := Positive[float64]()(0.1); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := Positive[float64]()(-0.1); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestNegative(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"negative value", -5, false},
		{"zero", 0, true},
		{"positive value", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Negative[int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Negative()(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNonZero(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"positive value", 5, false},
		{"negative value", -5, false},
		{"zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NonZero[int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NonZero()(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestMultipleOf(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		value   int
		wantErr bool
	}{
		{"10 is multiple of 5", 5, 10, false},
		{"15 is multiple of 5", 5, 15, false},
		{"7 is not multiple of 5", 5, 7, true},
		{"0 is multiple of any", 5, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MultipleOf(tt.n)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultipleOf(%d)(%d) error = %v, wantErr %v", tt.n, tt.value, err, tt.wantErr)
			}
		})
	}

	t.Run("panics with zero divisor", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got nil")
			}
		}()
		MultipleOf(0)
	})
}
