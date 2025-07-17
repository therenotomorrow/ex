package ex

import (
	"errors"
)

const (
	// ErrUnexpected represents an unexpected error.
	ErrUnexpected LError = "unexpected"
)

type (
	// LError is an error for const information.
	LError string

	// XError is an error for extended information.
	XError struct {
		label error
		cause error
	}
	// L is an LError alias.
	L = LError
	// X is an XError alias.
	X = XError
)

// New creates a new XError from a string.
func New(text string) *XError {
	return &XError{label: LError(text), cause: nil}
}

// From creates an XError from another error.
func From(other error) *XError {
	var xer *XError
	if errors.As(other, &xer) {
		return &XError{label: xer.label, cause: xer.cause}
	}

	return &XError{label: other, cause: nil}
}

// Unexpected creates an XError with ErrUnexpected as its error and a given cause.
func Unexpected(cause error) error {
	return &XError{label: ErrUnexpected, cause: cause}
}

// Must returns the value or panics if an error occurs.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(Cause(err))
	}

	return t
}

// MustDo executes panics if an error occurs.
func MustDo(err error) {
	if err != nil {
		panic(Cause(err))
	}
}

// Cause returns the root cause of an error.
func Cause(err error) error {
	for err != nil {
		var xer *XError
		if !errors.As(err, &xer) {
			break
		}

		if xer.cause == nil {
			return xer.label
		}

		err = xer.cause
	}

	return err
}

// Because wraps a LError with an additional cause.
func (c LError) Because(cause error) error {
	return &XError{label: c, cause: cause}
}

// Reason wraps a LError with a reason string as its cause.
func (c LError) Reason(text string) error {
	return &XError{label: c, cause: LError(text)}
}

// Error implements the error interface for LError.
func (c LError) Error() string {
	return string(c)
}

// Because wraps an XError with an additional cause.
func (e *XError) Because(cause error) error {
	return &XError{label: e.label, cause: cause}
}

// Reason wraps an XError with a reason string as its cause.
func (e *XError) Reason(text string) error {
	return &XError{label: e.label, cause: LError(text)}
}

// Error implements the error interface for XError.
func (e *XError) Error() string {
	return e.label.Error()
}

// Unwrap returns the wrapped error in XError.
// Implements anonymous Unwrap interface to be compliant with errors.Is.
func (e *XError) Unwrap() error {
	return e.label
}

// Is checks if the XError or its cause matches a target error.
// Implements anonymous Is interface to be compliant with errors.Is.
func (e *XError) Is(target error) bool {
	if errors.Is(e.label, target) {
		return true
	}

	if e.cause != nil {
		return errors.Is(e.cause, target)
	}

	return false
}
