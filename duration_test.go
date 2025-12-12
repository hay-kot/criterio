package criterio

import (
	"testing"
	"time"
)

func TestDurMin(t *testing.T) {
	tests := []struct {
		name    string
		min     time.Duration
		value   time.Duration
		wantErr bool
	}{
		{"duration equals min", time.Minute, time.Minute, false},
		{"duration above min", time.Minute, time.Hour, false},
		{"duration below min", time.Minute, time.Second, true},
		{"zero duration", time.Second, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DurMin(tt.min)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurMin(%v)(%v) error = %v, wantErr %v", tt.min, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestDurMax(t *testing.T) {
	tests := []struct {
		name    string
		max     time.Duration
		value   time.Duration
		wantErr bool
	}{
		{"duration equals max", time.Hour, time.Hour, false},
		{"duration below max", time.Hour, time.Minute, false},
		{"duration above max", time.Minute, time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DurMax(tt.max)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurMax(%v)(%v) error = %v, wantErr %v", tt.max, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestDurBetween(t *testing.T) {
	tests := []struct {
		name    string
		low     time.Duration
		high    time.Duration
		value   time.Duration
		wantErr bool
	}{
		{"duration equals low", time.Second, time.Hour, time.Second, false},
		{"duration equals high", time.Second, time.Hour, time.Hour, false},
		{"duration in range", time.Second, time.Hour, time.Minute, false},
		{"duration below range", time.Minute, time.Hour, time.Second, true},
		{"duration above range", time.Second, time.Minute, time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DurBetween(tt.low, tt.high)(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurBetween(%v, %v)(%v) error = %v, wantErr %v", tt.low, tt.high, tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestDurPositive(t *testing.T) {
	tests := []struct {
		name    string
		value   time.Duration
		wantErr bool
	}{
		{"positive duration", time.Second, false},
		{"zero duration", 0, true},
		{"negative duration", -time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DurPositive()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurPositive()(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
