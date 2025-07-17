// Package ex shows another way for working with the errors in Go.
// The idea is simple - use some "label" or "slug" for errors and
// wrap or unwrap the exact "cause" or "reason" within.
//
// Two main things:
//   - LError - create error as the constants (aka "label", "slug", etc.)
//   - XError - extend errors with the internals (aka "cause", "reason", etc.)
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
//			ErrUserNotFound = ex.LError("user not found")
//			ErrDatabase     = ex.LError("database error")
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
//		// Full error: user not found
//	}
package ex
