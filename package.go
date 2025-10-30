// Package ex provides a flexible approach to error handling in Go. It allows for
// creating chainable errors that combine a static, constant "identity" with a
// dynamic, underlying "cause".
//
// The core idea is to separate what an error is (its identity, e.g., "user not found")
// from why it happened (its cause, e.g., a specific database connection error).
//
// This is achieved with two main components:
//   - Error: A constant string type for defining sentinel error identities.
//   - XError: An interface for errors that can be chained.
//
// This pattern allows you to check for high-level errors using errors.Is while
// preserving the full, detailed context of the original problem, making debugging
// and logging much more effective.
//
// Best practices:
// - Define sentinel errors as package-level constants
// - Keep error chains shallow when possible (2-3 levels)
// - Use Error.Reason for simple text descriptions
// - Use Error.Because when wrapping existing errors
//
// # Example
//
// This example demonstrates how to define domain-specific errors and wrap them
// to create a rich error chain.
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
//		const (
//			ErrUserNotFound = ex.Error("user not found")
//			ErrDatabase     = ex.Error("database error")
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
//			// 6. You can also check for intermediate errors in the chain.
//			if errors.Is(err, ErrDatabase) {
//				fmt.Println("The cause was a database error.")
//			}
//
//			// 7. The full error message shows the complete chain of identities.
//			fmt.Printf("Full error: %s\n", err)
//		}
//	}
//
//	// Output:
//	// Error is user not found.
//	// The cause was a database error.
//	// Full error: user not found: database error: connection reset by peer
package ex
