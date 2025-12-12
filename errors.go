// Package criterio provides utilities for marking and handling validation errors.
package criterio

import (
	"fmt"
	"strings"
)

// FieldError represents a validation error for a specific field. Generally, you won't deal with
// this type directly and instead you'll want to use [FieldErrorsBuilder] to build errors.
type FieldError struct {
	Field   string
	Message string
}

// Error returns the error message for this field error.
func (e FieldError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// FieldErrorsBuilder is used to build a collection of field validation errors.
// It does not implement the error interface to prevent accidental returns without calling ToError().
type FieldErrorsBuilder []FieldError

// Append adds a field error to the collection.
// Returns the updated FieldErrorsBuilder for convenient chaining.
//
// Example usage:
//
//	var errs validation.FieldErrorsBuilder
//	errs = errs.Append("name", "cannot be empty")
//	errs = errs.Append("duration", "must be positive")
//	return errs.ToError()
func (b FieldErrorsBuilder) Append(field, message string) FieldErrorsBuilder {
	return append(b, FieldError{Field: field, Message: message})
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
		// In practice, this should never happen, but we include it as a case to catch any
		// incorrect usage in the logs.
		return "validation failed with zero error collected"
	}
	if len(e) == 1 {
		return e[0].Error()
	}

	bldr := &strings.Builder{}
	bldr.WriteString("validation failed: ")
	for i, err := range e {
		if i > 0 {
			bldr.WriteRune(';')
			bldr.WriteRune(' ')
		}

		bldr.WriteString(err.Error())
	}

	return bldr.String()
}

// NewFieldErrors creates a new FieldErrors from a single field error.
func NewFieldErrors(field, message string) FieldErrors {
	return FieldErrors{FieldError{Field: field, Message: message}}
}
