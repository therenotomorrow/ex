package ex_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/therenotomorrow/ex"
)

// Demonstrates creating a new, basic error from a string.
func ExampleNew() {
	err := ex.New("repository: user not found")

	fmt.Println(err)
	// Output:
	// repository: user not found
}

// Demonstrates how to convert a standard library or third-party error
// into an XError, allowing it to be part of a chain while preserving its original identity.
func ExampleCast() {
	var (
		// Simulate an error from an external package.
		originalErr = io.EOF
		// Cast the standard error to an XError.
		err = ex.Cast(originalErr)
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

// Demonstrates how to assert that a function call must not return an error.
// If an error is present, Must panics. This is useful for setup code where an
// error is unrecoverable and should halt execution immediately.
func ExampleMust() {
	// Case 1: The function succeeds and returns a value.
	value := ex.Must("critical data", nil)

	fmt.Println(value)

	// Case 2: The function returns an error and Must panics.
	// We use a deferred recover to gracefully handle the panic for this example.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	ex.Must("some value", errors.New("something went wrong"))
	// Output:
	// critical data
	// Recovered from panic: something went wrong
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
