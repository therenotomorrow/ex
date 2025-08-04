package ex_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/therenotomorrow/ex"
)

// Demonstrates creating a basic error.
func ExampleNew() {
	err := ex.New("repository: user not found")
	fmt.Println(err)
	// Output:
	// repository: user not found
}

// Demonstrates converting a standard error into an ExtraError.
func ExampleCast() {
	// Simulate an error from an external package
	originalErr := io.EOF

	// Wrap it in an ExtraError to add context
	err := ex.Cast(originalErr)
	fmt.Println(err)

	// You can still check for the original error type
	if errors.Is(err, io.EOF) {
		fmt.Println("Error is io.EOF")
	}
	// Output:
	// EOF
	// Error is io.EOF
}

// Demonstrates wrapping an unexpected error with a standard message.
func ExampleUnexpected() {
	// An unexpected database error occurs
	dbErr := errors.New("connection refused")

	// Wrap it as a critical, unexpected error
	err := ex.Unexpected(dbErr)
	fmt.Println(err)

	// The main error message is standardized
	if errors.Is(err, ex.ErrUnexpected) {
		fmt.Println("This was an unexpected error.")
	}
	// Output:
	// unexpected: connection refused
	// This was an unexpected error.
}

// Demonstrates wrapping a critical error with a standard message.
func ExampleCritical() {
	// An unexpected database error occurs
	dbErr := errors.New("connection refused")

	// Wrap it as a critical, unexpected error
	err := ex.Critical(dbErr)
	fmt.Println(err)

	// The main error message is standardized
	if errors.Is(err, ex.ErrCritical) {
		fmt.Println("This was a critical error.")
	}
	// Output:
	// critical: connection refused
	// This was a critical error.
}

// Demonstrates its use in situations where an error is not want.
func ExampleMust() {
	// This function simulates a call that should not fail
	mightReturnValue := func() (string, error) {
		return "critical data", nil
	}

	// This function simulates a call that will fail
	mightPanic := func() (string, error) {
		return "", errors.New("something went wrong")
	}

	// The successful case:
	value := ex.Must(mightReturnValue())
	fmt.Println(value)

	// The panic case:
	// We use a deferred recover to demonstrate the panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	ex.Must(mightPanic())
	// Output:
	// critical data
	// Recovered from panic: something went wrong
}

// Shows how to add a causal error for ConstError.
func ExampleError_Because() {
	ErrPayment := ex.Error("payment failed")
	apiErr := errors.New("stripe: invalid API key")

	err := ErrPayment.Because(apiErr)

	fmt.Println(err)
	// Output:
	// payment failed: stripe: invalid API key
}

// Shows how to add a cause with a simple text description for ConstError.
func ExampleError_Reason() {
	ErrValidation := ex.Error("validation failed")

	err := ErrValidation.Reason("email address is missing")

	fmt.Println(err)
	// Output:
	// validation failed: email address is missing
}
