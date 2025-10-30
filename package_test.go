package ex_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/therenotomorrow/ex"
)

// Demonstrates how to convert a standard library or third-party error
// into an XError, allowing it to be part of a chain while preserving its original identity.
func ExampleNew() {
	var (
		// Simulate an error from an external package.
		originalErr = io.EOF
		// Cast the standard error to an XError.
		err = ex.New(originalErr)
	)

	fmt.Println(err)

	// You can still check for the original error identity using errors.Is.
	if errors.Is(err, io.EOF) {
		fmt.Println("Error is io.EOF")
	}
	// Output:
	// EOF
	// Error is io.EOF
}

// Demonstrates how Panic is works.
// If an error is present, Panic panics with the ErrCritical. This is useful for setup code where an
// error is unrecoverable and should halt execution immediately.
func ExamplePanic() {
	// Case 1: The function succeeds and returns a value.
	ex.Panic(nil)

	// Case 2: The function returns an error and Panic panics.
	// We use a deferred recover to gracefully handle the panic for this example.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	ex.Panic(errors.New("something went wrong"))
	// Output:
	// Recovered from panic: critical: something went wrong
}

// Demonstrates how Skip is works.
// If an error is present, Skip ignores it and marks for code readers that this is skipped.
func ExampleSkip() {
	// Case 1: All fine with nil.
	ex.Skip(nil)

	// Case 2: From somewhere we received the error, just skip it.
	ex.Skip(errors.New("something went wrong"))

	// Output:
}

// Demonstrates how WithPanic is works.
// If an error is present, WithPanic panics with the ErrCritical.
// Same as Panic but more useful for return values.
func ExampleWithPanic() {
	// Case 1: The function succeeds and returns a value.
	got := ex.WithPanic("result", nil)
	fmt.Println("Our success result:", got)

	// Case 2: The function returns an error and WithPanic panics.
	// We use a deferred recover to gracefully handle the panic for this example.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	_ = ex.WithPanic("missing", errors.New("something went wrong"))
	// Output:
	// Our success result: result
	// Recovered from panic: critical: something went wrong
}

// Demonstrates how WithSkip is works.
// Same as Panic but more useful for return values.
func ExampleWithSkip() {
	// Case 1: All fine with nil.
	got := ex.WithSkip("kekes", nil)
	fmt.Println("Our first success:", got)

	// Case 2: From somewhere we received the error, just skip it.
	got = ex.WithSkip("memes", errors.New("something went wrong"))
	fmt.Println("Our second success:", got)

	// Output:
	// Our first success: kekes
	// Our second success: memes
}

// Shows how to wrap an error with the standard ErrUnexpected identity.
// This is useful for flagging errors that likely indicate a bug.
func ExampleUnexpected() {
	var (
		// Simulate an unexpected database error.
		dbErr = errors.New("connection refused")
		// Wrap it as an unexpected error.
		err = ex.Unexpected(dbErr)
	)

	fmt.Println(err)

	// You can now check for the generic "unexpected" error type.
	if errors.Is(err, ex.ErrUnexpected) {
		fmt.Println("This was an unexpected error.")
	}
	// Output:
	// unexpected: connection refused
	// This was an unexpected error.
}

// Shows how to wrap an error with the standard ErrCritical identity,
// signaling a severe, non-recoverable problem.
func ExampleCritical() {
	var (
		// Simulate a critical filesystem error.
		fsErr = errors.New("disk is full")
		// Wrap it as a critical error.
		err = ex.Critical(fsErr)
	)

	fmt.Println(err)

	// Check for the generic "critical" error type.
	if errors.Is(err, ex.ErrCritical) {
		fmt.Println("This was a critical error.")
	}
	// Output:
	// critical: disk is full
	// This was a critical error.
}

// Shows how to wrap an error with the standard ErrDummy identity.
// This is useful for testing and not meaningful errors.
func ExampleDummy() {
	var (
		// Simulate a dummy eol error.
		eolErr = errors.New("finite task")
		// Wrap it as a dummy error.
		err = ex.Dummy(eolErr)
	)

	fmt.Println(err)

	// Check for the generic "dummy" error type.
	if errors.Is(err, ex.ErrDummy) {
		fmt.Println("This was a dummy error.")
	}
	// Output:
	// dummy: finite task
	// This was a dummy error.
}

// Shows how to add a causal error to a sentinel Error, creating a chain of errors.
func ExampleError_Because() {
	// Define a sentinel error for your domain.
	const ErrPayment ex.Error = "payment failed"

	var (
		// Simulate a low-level API error.
		apiErr = errors.New("stripe: invalid API key")
		// Chain the errors together. ErrPayment is the identity, apiErr is the cause.
		err = ErrPayment.Because(apiErr)
	)

	fmt.Println(err)
	// Output:
	// payment failed: stripe: invalid API key
}

// Shows how to add a simple text description as the cause for a sentinel Error.
func ExampleError_Reason() {
	// Define a sentinel error for your domain.
	const ErrValidation ex.Error = "validation failed"

	// Add a specific, human-readable reason.
	err := ErrValidation.Reason("email address is missing")

	fmt.Println(err)
	// Output:
	// validation failed: email address is missing
}
