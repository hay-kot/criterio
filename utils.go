package criterio

import (
	"errors"
	"strconv"
)

// Static error messages for utility validators.
var errRequired = errors.New("is required")

// Validator validates a value and returns an error if validation fails.
type Validator[T any] func(val T) error

// Run executes validators against the value, stopping on the first failure.
// Returns nil if all validators pass, or a FieldErrors with the first failure.
func Run[T any](field string, val T, validators ...Validator[T]) error {
	for _, v := range validators {
		if err := v(val); err != nil {
			return NewFieldErrors(field, err)
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
			errs = errs.Append(field, err)
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
		return errRequired
	}
	return nil
}

// ValidateStruct combines multiple field validation results into a single error.
// Pass the results of Run/RunAll calls for each field.
//
// Non-FieldErrors errors (e.g., from external libraries) are included with an
// empty field name. For proper field context, wrap external errors with Nest
// or use Run/RunAll.
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
			errs = errs.Append("", err)
		}
	}
	return errs.ToError()
}

// Nest prefixes all field errors with a parent field name using dot notation.
// Useful for nested struct validation to create paths like "address.street".
//
// Non-FieldErrors errors are wrapped with the field name as context.
func Nest(field string, err error) error {
	if err == nil {
		return nil
	}
	var fieldErrs FieldErrors
	ok := errors.As(err, &fieldErrs)
	if !ok {
		return NewFieldErrors(field, err)
	}
	nested := make(FieldErrors, len(fieldErrs))
	for i, fe := range fieldErrs {
		if fe.Field == "" {
			nested[i] = FieldError{Field: field, Err: fe.Err}
		} else {
			nested[i] = FieldError{Field: field + "." + fe.Field, Err: fe.Err}
		}
	}
	return nested
}

// When returns a validator that only runs if the condition is true.
// If the condition is false, validation passes without running validators.
// Panics if no validators are provided.
func When[T any](condition bool, validators ...Validator[T]) Validator[T] {
	if len(validators) == 0 {
		panic("When: at least one validator is required")
	}
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
// Panics if no validators are provided.
func Or[T any](validators ...Validator[T]) Validator[T] {
	if len(validators) == 0 {
		panic("Or: at least one validator is required")
	}
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

// ValidateSlice validates each element in a slice using the provided function.
// Errors are prefixed with bracket notation (e.g., "items[0].name: is required").
func ValidateSlice[T any](field string, items []T, validate func(T) error) error {
	var errs FieldErrorsBuilder
	for i, item := range items {
		if err := validate(item); err != nil {
			prefix := field + "[" + strconv.Itoa(i) + "]"
			var fieldErrs FieldErrors
			if errors.As(err, &fieldErrs) {
				for _, fe := range fieldErrs {
					if fe.Field == "" {
						errs = errs.Append(prefix, fe.Err)
					} else {
						errs = errs.Append(prefix+"."+fe.Field, fe.Err)
					}
				}
			} else {
				errs = errs.Append(prefix, err)
			}
		}
	}
	return errs.ToError()
}
