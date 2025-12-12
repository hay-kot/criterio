package criterio

import (
	"testing"
	"time"
)

func TestTimeFuture(t *testing.T) {
	t.Run("future time passes", func(t *testing.T) {
		future := time.Now().Add(time.Hour)
		if err := TimeFuture()(future); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("past time fails", func(t *testing.T) {
		past := time.Now().Add(-time.Hour)
		if err := TimeFuture()(past); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("now fails", func(t *testing.T) {
		// time.Now() at check time will be equal or past
		now := time.Now()
		if err := TimeFuture()(now); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTimePast(t *testing.T) {
	t.Run("past time passes", func(t *testing.T) {
		past := time.Now().Add(-time.Hour)
		if err := TimePast()(past); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("future time fails", func(t *testing.T) {
		future := time.Now().Add(time.Hour)
		if err := TimePast()(future); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTimeAfter(t *testing.T) {
	reference := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"time after reference", reference.Add(time.Hour), false},
		{"time before reference", reference.Add(-time.Hour), true},
		{"time equals reference", reference, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TimeAfter(reference)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeAfter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeBefore(t *testing.T) {
	reference := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"time before reference", reference.Add(-time.Hour), false},
		{"time after reference", reference.Add(time.Hour), true},
		{"time equals reference", reference, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TimeBefore(reference)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeBefore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeBetween(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"time at start", start, false},
		{"time at end", end, false},
		{"time in range", time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), false},
		{"time before range", time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC), true},
		{"time after range", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TimeBetween(start, end)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeBetween() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
