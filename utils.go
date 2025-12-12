package criterio

import "fmt"

// Validator validates a value and returns an error if validation fails.
type Validator[T any] func(val T) error

// Run executes all validators against the value and collects any errors.
// Returns nil if all validators pass, or a FieldErrors containing all failures.
func Run[T any](field string, val T, validators ...Validator[T]) error {
	var errs FieldErrorsBuilder
	for _, v := range validators {
		if err := v(val); err != nil {
			errs = errs.Append(field, err.Error())
		}
	}
	return errs.ToError()
}

// New creates a reusable validator function for a specific field.
// Returns a function that validates values and returns nil or FieldErrors.
func New[T any](field string, validators ...Validator[T]) func(val T) error {
	return func(val T) error {
		var errs FieldErrorsBuilder
		for _, v := range validators {
			if err := v(val); err != nil {
				errs = errs.Append(field, err.Error())
			}
		}
		return errs.ToError()
	}
}

// Required returns a validator that checks if a value is non-zero.
func Required[T comparable]() Validator[T] {
	return func(val T) error {
		var zero T
		if val == zero {
			return fmt.Errorf("is required")
		}
		return nil
	}
}
