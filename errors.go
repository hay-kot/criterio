// Package criterio provides utilities for marking and handling validation errors.
package criterio

import (
	"strings"
)

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field string
	Err   error
}

// Error returns the error message for this field error.
func (e FieldError) Error() string {
	if e.Field == "" {
		return e.Err.Error()
	}
	// Use string concatenation - faster than fmt.Sprintf for simple cases
	return e.Field + ": " + e.Err.Error()
}

// Unwrap returns the underlying error for use with errors.Is and errors.As.
func (e FieldError) Unwrap() error {
	return e.Err
}

// FieldErrorsBuilder is used to build a collection of field validation errors.
// It does not implement the error interface to prevent accidental returns without calling ToError().
type FieldErrorsBuilder []FieldError

// Append adds a field error to the collection.
// Returns the updated FieldErrorsBuilder for convenient chaining.
//
// Example usage:
//
//	var errs criterio.FieldErrorsBuilder
//	errs = errs.Append("name", errRequired)
//	errs = errs.Append("duration", errPositive)
//	return errs.ToError()
func (b FieldErrorsBuilder) Append(field string, err error) FieldErrorsBuilder {
	return append(b, FieldError{Field: field, Err: err})
}

// ToError converts FieldErrorsBuilder to an error.
// Returns nil if there are no errors, preventing the empty-slice-as-non-nil-error issue.
func (b FieldErrorsBuilder) ToError() error {
	if len(b) == 0 {
		return nil
	}
	return FieldErrors(b)
}

// FieldErrors is a collection of field validation errors that implements the error interface.
type FieldErrors []FieldError

// Error returns a human-readable error message combining all field errors.
func (e FieldErrors) Error() string {
	if len(e) == 0 {
		return "validation failed with zero error collected"
	}
	if len(e) == 1 {
		return e[0].Error()
	}

	bldr := &strings.Builder{}
	bldr.WriteString("validation failed: ")
	for i, err := range e {
		if i > 0 {
			bldr.WriteString("; ")
		}
		bldr.WriteString(err.Error())
	}

	return bldr.String()
}

// NewFieldError creates a single FieldError.
func NewFieldError(field string, err error) FieldError {
	return FieldError{Field: field, Err: err}
}

// NewFieldErrors creates a new FieldErrors from a single field error.
func NewFieldErrors(field string, err error) FieldErrors {
	return FieldErrors{{Field: field, Err: err}}
}
