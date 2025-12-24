package criterio

import (
	"errors"
	"fmt"
)

// Validator validates a value and returns an error if validation fails.
type Validator[T any] func(val T) error

// Run executes validators against the value, stopping on the first failure.
// Returns nil if all validators pass, or a FieldErrors with the first failure.
func Run[T any](field string, val T, validators ...Validator[T]) error {
	for _, v := range validators {
		if err := v(val); err != nil {
			return NewFieldErrors(field, err.Error())
		}
	}
	return nil
}

// RunAll executes all validators against the value, collecting all failures.
// Returns nil if all validators pass, or a FieldErrors containing all failures.
func RunAll[T any](field string, val T, validators ...Validator[T]) error {
	var errs FieldErrorsBuilder
	for _, v := range validators {
		if err := v(val); err != nil {
			errs = errs.Append(field, err.Error())
		}
	}
	return errs.ToError()
}

// New creates a reusable validator function for a specific field.
// Stops on the first validation failure.
func New[T any](field string, validators ...Validator[T]) func(val T) error {
	return func(val T) error {
		return Run(field, val, validators...)
	}
}

// NewAll creates a reusable validator function for a specific field.
// Collects all validation failures.
func NewAll[T any](field string, validators ...Validator[T]) func(val T) error {
	return func(val T) error {
		return RunAll(field, val, validators...)
	}
}

// Required validates that a value is non-zero.
func Required[T comparable](val T) error {
	var zero T
	if val == zero {
		return fmt.Errorf("is required")
	}
	return nil
}

// ValidateStruct combines multiple field validation results into a single error.
// Pass the results of Run/RunAll calls for each field.
func ValidateStruct(validations ...error) error {
	var errs FieldErrorsBuilder
	for _, err := range validations {
		if err == nil {
			continue
		}
		var fieldErrs FieldErrors
		if errors.As(err, &fieldErrs) {
			errs = append(errs, fieldErrs...)
		} else {
			errs = errs.Append("", err.Error())
		}
	}
	return errs.ToError()
}

// Nest prefixes all field errors with a parent field name using dot notation.
// Useful for nested struct validation to create paths like "address.street".
func Nest(field string, err error) error {
	if err == nil {
		return nil
	}
	var fieldErrs FieldErrors
	ok := errors.As(err, &fieldErrs)
	if !ok {
		return NewFieldErrors(field, err.Error())
	}
	nested := make(FieldErrors, len(fieldErrs))
	for i, fe := range fieldErrs {
		if fe.Field == "" {
			nested[i] = FieldError{Field: field, Message: fe.Message}
		} else {
			nested[i] = FieldError{Field: field + "." + fe.Field, Message: fe.Message}
		}
	}
	return nested
}

// When returns a validator that only runs if the condition is true.
// If the condition is false, validation passes without running validators.
func When[T any](condition bool, validators ...Validator[T]) Validator[T] {
	return func(val T) error {
		if !condition {
			return nil
		}
		for _, v := range validators {
			if err := v(val); err != nil {
				return err
			}
		}
		return nil
	}
}

// SkipIf returns a validator that skips validation if the condition is true.
// If the condition is true, validation passes without running validators.
func SkipIf[T any](condition bool, validators ...Validator[T]) Validator[T] {
	return When(!condition, validators...)
}

// Or returns a validator that passes if any of the validators pass.
// Returns the last error if all validators fail.
func Or[T any](validators ...Validator[T]) Validator[T] {
	return func(val T) error {
		var lastErr error
		for _, v := range validators {
			if err := v(val); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}
		return lastErr
	}
}

// Not returns a validator that inverts the result of another validator.
// Passes if the wrapped validator fails, fails if the wrapped validator passes.
func Not[T any](v Validator[T], msg string) Validator[T] {
	return func(val T) error {
		if err := v(val); err == nil {
			return errors.New(msg)
		}
		return nil
	}
}
