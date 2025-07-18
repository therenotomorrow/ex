package ex

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// ErrUnexpected represents an unexpected error.
	ErrUnexpected ConstError = "unexpected"
)

type (
	// ConstError is an error for const information.
	ConstError string

	// ExtraError is an error for extended information.
	ExtraError struct {
		err   error
		cause error
	}
)

// C is an ConstError alias.
type C = ConstError

// New creates a new ExtraError from a string.
func New(text string) *ExtraError {
	return &ExtraError{err: ConstError(text), cause: nil}
}

// From creates an ExtraError from another error.
func From(other error) *ExtraError {
	var xer *ExtraError
	if errors.As(other, &xer) {
		return &ExtraError{err: xer.err, cause: xer.cause}
	}

	return &ExtraError{err: other, cause: nil}
}

// Unexpected creates an ExtraError with ErrUnexpected as its error and a given cause.
func Unexpected(cause error) error {
	return &ExtraError{err: ErrUnexpected, cause: cause}
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
		var xer *ExtraError
		if !errors.As(err, &xer) {
			break
		}

		if xer.cause == nil {
			return xer.err
		}

		err = xer.cause
	}

	return err
}

// Because wraps a ConstError with an additional cause.
func (c ConstError) Because(cause error) error {
	return &ExtraError{err: c, cause: cause}
}

// Reason wraps a ConstError with a reason string as its cause.
func (c ConstError) Reason(text string) error {
	return &ExtraError{err: c, cause: ConstError(text)}
}

// Error implements the error interface for ConstError.
func (c ConstError) Error() string {
	return string(c)
}

// String implements the fmt.Stringer interface for ConstError.
func (c ConstError) String() string {
	return string(c)
}

// Because wraps an ExtraError with an additional cause.
func (e *ExtraError) Because(cause error) error {
	return &ExtraError{err: e.err, cause: cause}
}

// Reason wraps an ExtraError with a reason string as its cause.
func (e *ExtraError) Reason(text string) error {
	return &ExtraError{err: e.err, cause: ConstError(text)}
}

// Error implements the error interface for ExtraError.
func (e *ExtraError) Error() string {
	if e.err == nil {
		return ""
	}

	msg := e.err.Error()

	cause := Cause(e.cause)
	if cause != nil {
		msg += fmt.Sprintf(" (%s)", cause)
	}

	return msg
}

// String implements the fmt.Stringer interface for ExtraError.
func (e *ExtraError) String() string {
	val := make(map[string]string)

	if e.err != nil {
		val["error"] = e.err.Error()
	}

	if e.cause != nil {
		val["cause"] = e.cause.Error()
	}

	bytes, _ := json.Marshal(val) //nolint:errchkjson // impossible error because of map

	return string(bytes)
}

// Unwrap returns the wrapped error in ExtraError.
// Implements anonymous Unwrap interface to be compliant with errors.Is.
func (e *ExtraError) Unwrap() error {
	return e.err
}

// Is checks if the ExtraError or its cause matches a target error.
// Implements anonymous Is interface to be compliant with errors.Is.
func (e *ExtraError) Is(target error) bool {
	if errors.Is(e.err, target) {
		return true
	}

	if e.cause != nil {
		return errors.Is(e.cause, target)
	}

	return false
}
