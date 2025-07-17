package ex_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/therenotomorrow/ex"
)

// ExampleNew demonstrates creating a basic error.
func ExampleNew() {
	err := ex.New("repository: user not found")
	fmt.Println(err)
	// Output:
	// repository: user not found
}

// ExampleFrom demonstrates converting a standard error into an XError.
func ExampleFrom() {
	// Simulate an error from an external package
	originalErr := io.EOF

	// Wrap it in an XError to add context
	err := ex.From(originalErr)
	fmt.Println(err)

	// You can still check for the original error type
	if errors.Is(err, io.EOF) {
		fmt.Println("Error is io.EOF")
	}
	// Output:
	// EOF
	// Error is io.EOF
}

// ExampleUnexpected demonstrates wrapping an unexpected error with a standard message.
func ExampleUnexpected() {
	// An unexpected database error occurs
	dbErr := errors.New("connection refused")

	// Wrap it as a critical, unexpected error
	err := ex.Unexpected(dbErr)
	fmt.Println(err)

	// We can find the root cause
	cause := ex.Cause(err)
	fmt.Println(cause)

	// The main error message is standardized
	if errors.Is(err, ex.ErrUnexpected) {
		fmt.Println("This was an unexpected error.")
	}
	// Output:
	// unexpected
	// connection refused
	// This was an unexpected error.
}

// ExampleMust demonstrates its use in situations where an error is not want.
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

// ExampleMustDo demonstrates its use in situations where an error is not want.
func ExampleMustDo() {
	// This function simulates a call that should not fail
	mightNoPanic := func() error {
		return nil
	}

	// This function simulates a call that will fail
	mightPanic := func() error {
		return errors.New("something went wrong")
	}

	// The successful case:
	ex.MustDo(mightNoPanic())

	// The panic case:
	// We use a deferred recover to demonstrate the panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	ex.MustDo(mightPanic())
	// Output:
	// Recovered from panic: something went wrong
}

// ExampleCause shows how to find the root cause of a nested error.
func ExampleCause() {
	// Create a deeply nested error
	rootCause := errors.New("root cause: disk is full")
	err := ex.New("service layer error").Because(
		ex.New("repository layer error").Because(rootCause),
	)

	// Extract the original, underlying error
	cause := ex.Cause(err)
	fmt.Println(cause)

	// Create another deeply nested error
	ioErr := errors.New("disk write error")
	dbErr := ex.New("failed to save user").Because(ioErr)
	authErr := ex.LError("permission denied").Because(dbErr)
	finalErr := ex.Unexpected(authErr)

	// Extract everything that we could
	fmt.Printf("Error message: %s\n", finalErr)
	fmt.Printf("Is critical: %t\n", errors.Is(finalErr, ex.ErrUnexpected))
	fmt.Printf("Is permission error: %t\n", errors.Is(finalErr, ex.LError("permission denied")))

	// With the another root cause
	rootCause2 := ex.Cause(finalErr)
	fmt.Printf("Root cause: %s\n", rootCause2)
	// Output:
	// root cause: disk is full
	// Error message: unexpected
	// Is critical: true
	// Is permission error: true
	// Root cause: disk write error
}

// ExampleLError_Because shows how to add a causal error.
func ExampleLError_Because() {
	ErrPayment := ex.LError("payment failed")
	apiErr := errors.New("stripe: invalid API key")

	err := ErrPayment.Because(apiErr)

	fmt.Println(err)
	fmt.Println("Cause:", ex.Cause(err))
	// Output:
	// payment failed
	// Cause: stripe: invalid API key
}

// ExampleLError_Reason shows how to add a cause with a simple text description.
func ExampleLError_Reason() {
	ErrValidation := ex.LError("validation failed")

	err := ErrValidation.Reason("email address is missing")

	fmt.Println(err)
	fmt.Println("Cause:", ex.Cause(err))
	// Output:
	// validation failed
	// Cause: email address is missing
}

// ExampleXError_Because shows how to add a causal error.
func ExampleXError_Because() {
	ErrPayment := ex.New("payment failed")
	apiErr := errors.New("stripe: invalid API key")

	err := ErrPayment.Because(apiErr)

	fmt.Println(err)
	fmt.Println("Cause:", ex.Cause(err))
	// Output:
	// payment failed
	// Cause: stripe: invalid API key
}

// ExampleXError_Reason shows how to add a cause with a simple text description.
func ExampleXError_Reason() {
	ErrValidation := ex.New("validation failed")

	err := ErrValidation.Reason("email address is missing")

	fmt.Println(err)
	fmt.Println("Cause:", ex.Cause(err))
	// Output:
	// validation failed
	// Cause: email address is missing
}
