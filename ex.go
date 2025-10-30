package ex

import (
	"errors"
	"strings"
)

const (
	// ErrUnexpected represents an unexpected error, typically a bug or logic error.
	ErrUnexpected Error = "unexpected"

	// ErrCritical represents a critical, non-recoverable error.
	ErrCritical Error = "critical"

	// ErrDummy is used as a mock-like error for testing purposes.
	ErrDummy Error = "dummy"
)

var (
	_ XError = Error("")
	_ XError = (*xError)(nil)
)

// XError defines an interface for chainable errors.
// It allows for adding context and a causal chain to standard errors.
type XError interface {
	error
	// Reason adds a descriptive string as the cause of the error.
	Reason(text string) error
	// Because adds an existing error as the cause of the current error.
	Because(cause error) error
}

// New converts a standard error into an XError.
func New(err error) XError {
	if err == nil {
		return nil
	}

	var xer *xError
	if errors.As(err, &xer) {
		return &xError{error: xer.error, cause: xer.cause}
	}

	return &xError{error: err, cause: nil}
}

// Expose unwraps an error to reveal its internal components: the primary error and its cause.
// Panics with "invalid error type" if the error is not of type *xError.
func Expose(err error) (error, error) {
	var xer *xError
	if !errors.As(err, &xer) {
		panic("invalid error type")
	}

	return xer.error, xer.cause
}

// Panic panics if an error is present. Useful for handling critical situations that should halt execution.
func Panic(err error) {
	if err != nil {
		panic(Critical(err))
	}
}

// Skip marks the error as ignored or suppressed.
// Useful for deliberately ignoring errors instead
// of using default error handling mechanics.
func Skip(_ error) {}

// WithPanic does the same as Panic but returns the incoming value.
func WithPanic[T any](t T, err error) T {
	if err != nil {
		panic(Critical(err))
	}

	return t
}

// WithSkip does the same as Skip but returns the incoming value.
func WithSkip[T any](t T, _ error) T {
	return t
}

// Unexpected creates a new error with ErrUnexpected as the root and sets the cause.
// If the cause is nil, the result error will also be nil.
func Unexpected(cause error) error {
	if cause == nil {
		return nil
	}

	return &xError{error: ErrUnexpected, cause: cause}
}

// Critical creates a new error with ErrCritical as the root and sets the cause.
// If the cause is nil, the result error will also be nil.
func Critical(cause error) error {
	if cause == nil {
		return nil
	}

	return &xError{error: ErrCritical, cause: cause}
}

// Dummy creates a new error with ErrDummy as the root and sets the cause.
// If the cause is nil, the result error will also be nil.
func Dummy(cause error) error {
	if cause == nil {
		return nil
	}

	return &xError{error: ErrDummy, cause: cause}
}

// Error is a constant string-based error type.
type Error string

// Because creates a new xError, using the current Error as the root and setting the provided error as the cause.
func (c Error) Because(cause error) error {
	return &xError{error: c, cause: cause}
}

// Reason creates a new xError, using the current Error as the root and a new error from `text` as the cause.
func (c Error) Reason(text string) error {
	return &xError{error: c, cause: Error(text)}
}

// Error returns the string representation of the Error, satisfying the standard `error` interface.
func (c Error) Error() string {
	return string(c)
}

// xError is an implementation of XError that holds a primary error and a causal error.
// This structure allows for creating a chain of errors to provide rich context.
type xError struct {
	error error // The primary error identity.
	cause error // The underlying cause of the primary error (can be nil).
}

// Because creates a new xError, preserving the original primary error but replacing its cause.
func (e *xError) Because(cause error) error {
	return &xError{error: e.error, cause: cause}
}

// Reason creates a new xError, preserving the original primary error
// but replacing its cause with a new error from `text`.
func (e *xError) Reason(text string) error {
	return &xError{error: e.error, cause: Error(text)}
}

// Error flattens the error chain into a single, colon-separated string.
// It recursively traverses the cause chain to build the final error message.
func (e *xError) Error() string {
	var builder strings.Builder

	builder.WriteString(e.error.Error())

	for cause := e.cause; cause != nil; {
		var xer *xError
		if errors.As(cause, &xer) {
			if xer.error != nil {
				builder.WriteString(": ")
				builder.WriteString(xer.error.Error())
			}

			cause = xer.cause
		} else {
			builder.WriteString(": ")
			builder.WriteString(cause.Error())

			break
		}
	}

	return builder.String()
}

// Unwrap returns the primary error, allowing compatibility with `errors.As`.
// Note: It does not unwrap the `cause`. For that, see the `Is` method or `Expose` function.
func (e *xError) Unwrap() error {
	return e.error
}

// Is checks if the target error matches the primary error or any error in the cause chain.
// This makes xError fully compatible with `errors.Is`.
func (e *xError) Is(target error) bool {
	if errors.Is(e.error, target) {
		return true
	}

	if e.cause != nil {
		return errors.Is(e.cause, target)
	}

	return false
}
