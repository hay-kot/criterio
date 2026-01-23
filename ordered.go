package criterio

import (
	"cmp"
	"errors"
	"fmt"
)

// Static error messages for ordered validators.
var (
	errPositive = errors.New("must be positive")
	errNegative = errors.New("must be negative")
	errNonZero  = errors.New("must not be zero")
)

// Integer is a constraint for all integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// SignedNumber is a constraint for signed numeric types that can be negative.
type SignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Min returns a validator that checks if a value is at least min.
// Works with any ordered type (integers, floats, strings, etc.).
func Min[T cmp.Ordered](min T) Validator[T] {
	return func(val T) error {
		if val < min {
			return fmt.Errorf("must be at least %v", min)
		}
		return nil
	}
}

// Max returns a validator that checks if a value is at most max.
// Works with any ordered type (integers, floats, strings, etc.).
func Max[T cmp.Ordered](max T) Validator[T] {
	return func(val T) error {
		if val > max {
			return fmt.Errorf("must be at most %v", max)
		}
		return nil
	}
}

// Between returns a validator that checks if a value is between low and high (inclusive).
// Works with any ordered type (integers, floats, strings, etc.).
func Between[T cmp.Ordered](low, high T) Validator[T] {
	return func(val T) error {
		if val < low || val > high {
			return fmt.Errorf("must be between %v and %v", low, high)
		}
		return nil
	}
}

// Positive returns a validator that checks if a value is greater than zero.
// Works with signed numeric types (integers and floats).
func Positive[T SignedNumber]() Validator[T] {
	return func(val T) error {
		if val <= 0 {
			return errPositive
		}
		return nil
	}
}

// Negative returns a validator that checks if a value is less than zero.
// Works with signed numeric types (integers and floats).
func Negative[T SignedNumber]() Validator[T] {
	return func(val T) error {
		if val >= 0 {
			return errNegative
		}
		return nil
	}
}

// NonZero returns a validator that checks if a value is not zero.
// Works with signed numeric types (integers and floats).
func NonZero[T SignedNumber]() Validator[T] {
	return func(val T) error {
		if val == 0 {
			return errNonZero
		}
		return nil
	}
}

// MultipleOf returns a validator that checks if a value is divisible by n.
// Works with integer types only.
// Panics if n is zero.
func MultipleOf[T Integer](n T) Validator[T] {
	if n == 0 {
		panic("MultipleOf: divisor cannot be zero")
	}
	return func(val T) error {
		if val%n != 0 {
			return fmt.Errorf("must be a multiple of %v", n)
		}
		return nil
	}
}
