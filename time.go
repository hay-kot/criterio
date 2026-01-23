package criterio

import (
	"errors"
	"fmt"
	"time"
)

// Static error messages for time validators.
var (
	errTimeFuture = errors.New("must be in the future")
	errTimePast   = errors.New("must be in the past")
)

// TimeFuture returns a validator that checks if a time is in the future.
func TimeFuture() Validator[time.Time] {
	return func(val time.Time) error {
		if !val.After(time.Now()) {
			return errTimeFuture
		}
		return nil
	}
}

// TimePast returns a validator that checks if a time is in the past.
func TimePast() Validator[time.Time] {
	return func(val time.Time) error {
		if !val.Before(time.Now()) {
			return errTimePast
		}
		return nil
	}
}

// TimeAfter returns a validator that checks if a time is after a given time.
func TimeAfter(t time.Time) Validator[time.Time] {
	return func(val time.Time) error {
		if !val.After(t) {
			return fmt.Errorf("must be after %s", t.Format(time.RFC3339))
		}
		return nil
	}
}

// TimeBefore returns a validator that checks if a time is before a given time.
func TimeBefore(t time.Time) Validator[time.Time] {
	return func(val time.Time) error {
		if !val.Before(t) {
			return fmt.Errorf("must be before %s", t.Format(time.RFC3339))
		}
		return nil
	}
}

// TimeBetween returns a validator that checks if a time is between start and end (inclusive).
func TimeBetween(start, end time.Time) Validator[time.Time] {
	return func(val time.Time) error {
		if val.Before(start) || val.After(end) {
			return fmt.Errorf("must be between %s and %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
		}
		return nil
	}
}
