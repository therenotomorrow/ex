// Package ex shows another way for working with the errors in Go.
// The idea is simple - use some "err" or "slug" for errors and
// wrap or unwrap the exact "cause" or "reason" within.
//
// Two main things:
//   - ConstError - create error as the constants (aka "err", "slug", etc.)
//   - ExtraError - extend errors with the internals (aka "cause", "reason", etc.)
//
// Example will be better:
//
//	package main
//
//	import (
//		"errors"
//		"fmt"
//
//		"github.com/therenotomorrow/ex"
//	)
//
//	func main() {
//		// 1. Define sentinel errors for your domain. These are the error "identities".
//		var (
//			ErrUserNotFound = ex.ConstError("user not found")
//			ErrDatabase     = ex.ConstError("database error")
//		)
//
//		// 2. Simulate a low-level error (the root cause).
//		ioErr := errors.New("connection reset by peer")
//
//		// 3. Create a function that returns a wrapped error.
//		// The error's identity is `ErrUserNotFound`, but it preserves the full cause chain.
//		findUser := func() error {
//			// The database layer wraps the low-level I/O error.
//			dbErr := ErrDatabase.Because(ioErr)
//			// The service layer wraps the database error with a more specific identity.
//			return ErrUserNotFound.Because(dbErr)
//		}
//
//		// 4. Handle the error.
//		err := findUser()
//		if err != nil {
//			// 5. Check against the specific sentinel error using errors.Is.
//			// This works even though the error is deeply nested.
//			if errors.Is(err, ErrUserNotFound) {
//				fmt.Println("Error is user not found.")
//			}
//
//			// 6. You can also check for intermediate errors.
//			if errors.Is(err, ErrDatabase) {
//				fmt.Println("The cause was a database error.")
//			}
//
//			// 7. Extract the root cause for detailed logging or debugging.
//			rootCause := ex.Cause(err)
//			fmt.Printf("Root cause: %s\n", rootCause)
//
//			// 8. The default error message shows the top-level identity.
//			fmt.Printf("Full error: %s\n", err)
//		}
//
//		// Output:
//		// Error is user not found.
//		// The cause was a database error.
//		// Root cause: connection reset by peer
//		// Full error: user not found (connection reset by peer)
//	}
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

// C is a ConstError alias.
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

	text := e.err.Error()

	cause := Cause(e.cause)
	if cause != nil {
		text += fmt.Sprintf(" (%s)", cause)
	}

	return text
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

	bytes, _ := json.Marshal(val) //nolint:errchkjson // impossible error because of `map`

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

// Test uses for testing purposes. Also, you can use it to expose ExtraError internals.
func Test(err error) (error, error) {
	var xer *ExtraError
	if !errors.As(err, &xer) {
		panic("invalid error type")
	}

	return xer.err, xer.cause
}
