package criterio

import (
	"errors"
	"fmt"
	"time"
)

// Static error messages for duration validators.
var errDurPositive = errors.New("must be positive")

// DurMin returns a validator that checks if a duration is at least min.
func DurMin(min time.Duration) Validator[time.Duration] {
	return func(val time.Duration) error {
		if val < min {
			return fmt.Errorf("must be at least %s", min)
		}
		return nil
	}
}

// DurMax returns a validator that checks if a duration is at most max.
func DurMax(max time.Duration) Validator[time.Duration] {
	return func(val time.Duration) error {
		if val > max {
			return fmt.Errorf("must be at most %s", max)
		}
		return nil
	}
}

// DurBetween returns a validator that checks if a duration is between low and high (inclusive).
func DurBetween(low, high time.Duration) Validator[time.Duration] {
	return func(val time.Duration) error {
		if val < low || val > high {
			return fmt.Errorf("must be between %s and %s", low, high)
		}
		return nil
	}
}

// DurPositive validates that a duration is positive (greater than zero).
func DurPositive(val time.Duration) error {
	if val <= 0 {
		return errDurPositive
	}
	return nil
}
