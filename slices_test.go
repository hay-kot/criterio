package criterio

import (
	"strings"
	"testing"
)

func TestSliceLenMin(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		value   []int
		wantErr bool
	}{
		{"length equals min", 3, []int{1, 2, 3}, false},
		{"length above min", 3, []int{1, 2, 3, 4}, false},
		{"length below min", 3, []int{1, 2}, true},
		{"empty slice", 1, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SliceLenMin[int](tt.min)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceLenMin(%d)(%v) error = %v, wantErr %v", tt.min, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestSliceLenMax(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		value   []int
		wantErr bool
	}{
		{"length equals max", 3, []int{1, 2, 3}, false},
		{"length below max", 3, []int{1, 2}, false},
		{"length above max", 3, []int{1, 2, 3, 4}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SliceLenMax[int](tt.max)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceLenMax(%d)(%v) error = %v, wantErr %v", tt.max, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestSliceLenBetween(t *testing.T) {
	tests := []struct {
		name    string
		low     int
		high    int
		value   []int
		wantErr bool
	}{
		{"length equals low", 2, 4, []int{1, 2}, false},
		{"length equals high", 2, 4, []int{1, 2, 3, 4}, false},
		{"length in range", 2, 4, []int{1, 2, 3}, false},
		{"length below range", 2, 4, []int{1}, true},
		{"length above range", 2, 4, []int{1, 2, 3, 4, 5}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SliceLenBetween[int](tt.low, tt.high)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceLenBetween(%d, %d)(%v) error = %v, wantErr %v", tt.low, tt.high, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestSliceNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   []int
		wantErr bool
	}{
		{"non-empty slice", []int{1, 2, 3}, false},
		{"single element", []int{1}, false},
		{"empty slice", []int{}, true},
		{"nil slice", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SliceNotEmpty[int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceNotEmpty()(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestSliceUnique(t *testing.T) {
	tests := []struct {
		name    string
		value   []int
		wantErr bool
	}{
		{"all unique", []int{1, 2, 3}, false},
		{"with duplicates", []int{1, 2, 2, 3}, true},
		{"empty slice", []int{}, false},
		{"single element", []int{1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SliceUnique[int]()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceUnique()(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}

	t.Run("string slice", func(t *testing.T) {
		if err := SliceUnique[string]()([]string{"a", "b", "c"}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := SliceUnique[string]()([]string{"a", "b", "a"}); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSliceEach(t *testing.T) {
	t.Run("all elements pass", func(t *testing.T) {
		validator := SliceEach(Min(0))
		if err := validator([]int{1, 2, 3}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("element fails", func(t *testing.T) {
		validator := SliceEach(Min(0))
		err := validator([]int{1, -2, 3})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "element 1") {
			t.Errorf("expected error to mention element index, got: %v", err)
		}
	})

	t.Run("empty slice passes", func(t *testing.T) {
		validator := SliceEach(Min(0))
		if err := validator([]int{}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("with string validator", func(t *testing.T) {
		validator := SliceEach(StrNotEmpty())
		if err := validator([]string{"a", "b", "c"}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := validator([]string{"a", "", "c"}); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
