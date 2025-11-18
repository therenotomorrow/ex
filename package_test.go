package ex_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/therenotomorrow/ex"
)

// Demonstrates how to create an XError from the usual string.
func ExampleNew() {
	// Create the error.
	err := ex.New("your cool message")

	fmt.Println(err)

	// Now it should be wrapped as a const error identity using errors.Is.
	if errors.Is(err, ex.Error("your cool message")) {
		fmt.Println("Error is", err.Error())
	}
	// Output:
	// your cool message
	// Error is your cool message
}

// Demonstrates how to convert a standard library or third-party error
// into an XError, allowing it to be part of a chain while preserving its original identity.
func ExampleConv() {
	var (
		// Simulate an error from an external package.
		originalErr = io.EOF
		// Conv the standard error to an XError.
		err = ex.Conv(originalErr)
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

// Demonstrates how to use Expose to inspect error components.
// This is useful for logging, debugging, or custom error handling logic.
func ExampleExpose() {
	const ErrDatabase ex.Error = "database error"

	var (
		rootCause = errors.New("connection timeout")
		err       = ErrDatabase.Because(rootCause)
	)

	// Expose the error to get its components
	primary, cause := ex.Expose(err)

	fmt.Printf("Primary: %v\n", primary)
	fmt.Printf("Cause: %v\n", cause)

	// Expose works with the standard error as proxy (error, nil)
	primary, cause = ex.Expose(rootCause)
	fmt.Printf("Primary: %v\n", primary)
	fmt.Printf("Cause: %v\n", cause)
	// Output:
	// Primary: database error
	// Cause: connection timeout
	// Primary: connection timeout
	// Cause: <nil>
}

// Demonstrates how Panic works.
// If an error is present, Panic panics with the ErrCritical. This is useful for setup code where an
// error is unrecoverable and should halt execution immediately.
func ExamplePanic() {
	// Case 1: No error, so nothing happens.
	ex.Panic(nil)

	// Case 2: An error is passed and Panic panics.
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

// Demonstrates how Skip works.
// If an error is present, Skip ignores it and marks for code readers that this is skipped.
func ExampleSkip() {
	// Case 1: All fine with nil.
	ex.Skip(nil)

	// Case 2: From somewhere we received the error, just skip it.
	ex.Skip(errors.New("something went wrong"))

	// Output:
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

// Shows how to wrap an error with the standard ErrUnknown identity.
// This is useful for flagging errors that are likely unknown, unregistered, or unexpected issues.
func ExampleUnknown() {
	var (
		// Simulate an unknown error.
		wowErr = errors.New("wow we have an error")
		// Wrap it as an unknown error.
		err = ex.Unknown(wowErr)
	)

	fmt.Println(err)

	// You can now check for the generic "unknown" error type.
	if errors.Is(err, ex.ErrUnknown) {
		fmt.Println("This was an unknown error.")
	}
	// Output:
	// unknown: wow we have an error
	// This was an unknown error.
}

// Shows how to wrap an error with the standard ErrCritical identity,
// signaling a severe, non-recoverable problem. And yes - it panics.
func ExampleCritical() {
	// Case 1: No error, so nothing happens (returns result, no panic).
	res := ex.Critical("we got result", nil)
	fmt.Println(res)

	// Simulate a critical filesystem error.
	fsErr := errors.New("disk is full")

	// Case 2: An error is passed and Critical panics.
	// We use a deferred recover to gracefully handle the panic for this example.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Wrap it as a critical error and it will panic.
	_ = ex.Critical("omg", fsErr)
	// Output:
	// we got result
	// Recovered from panic: critical: disk is full
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
