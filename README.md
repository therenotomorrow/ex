# `ex`

Shows another way for working with the ~~ex~~ ( :dancer: ) errors in Go. Taste it! :heart:

<div>
  <a href="https://github.com/therenotomorrow/ex/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/therenotomorrow/ex?color=FBC02D" alt="GitHub releases">
  </a>
  <a href="https://go.dev/doc/go1.25" target="_blank">
    <img src="https://img.shields.io/badge/Go-%3E%3D%201.25-blue.svg" alt="Go 1.25">
  </a>
  <a href="https://pkg.go.dev/github.com/therenotomorrow/ex" target="_blank">
    <img src="https://godoc.org/github.com/therenotomorrow/ex?status.svg" alt="Go reference">
  </a>
  <a href="https://github.com/therenotomorrow/ex/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/therenotomorrow/ex?color=388E3C" alt="License">
  </a>
  <a href="https://github.com/therenotomorrow/ex/actions/workflows/ci.yml" target="_blank">
    <img src="https://github.com/therenotomorrow/ex/actions/workflows/ci.yml/badge.svg" alt="ci status">
  </a>
  <a href="https://goreportcard.com/report/github.com/therenotomorrow/ex" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/therenotomorrow/ex" alt="Go report">
  </a>
  <a href="https://codecov.io/gh/therenotomorrow/ex" target="_blank">
    <img src="https://img.shields.io/codecov/c/github/therenotomorrow/ex?color=546E7A" alt="Codecov">
  </a>
</div>

## Installation

```shell
go get github.com/therenotomorrow/ex@latest
```

Usage example:

```go
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/therenotomorrow/ex"
)

// 1. Define sentinel errors for your application domain.
// These act as stable identifiers for known error conditions.
const (
	ErrConfigValidation ex.Error = "config validation failed"
	ErrFileAccess       ex.Error = "file access error"
)

// 2. Create a function that simulates a real-world error chain.
func loadConfig(_ string) error {
	// Simulate a low-level OS error (the root cause).
	// For this example, we'll use a standard library error.
	underlyingErr := os.ErrPermission

	// The file access layer wraps the specific OS error with our domain error.
	accessErr := ErrFileAccess.Because(underlyingErr)

	// The business logic layer wraps the access error with a higher-level reason.
	businessErr := ErrConfigValidation.Reason("user section is missing")

	// Chain multiple errors: set accessErr as the cause of businessErr
	return ex.Conv(businessErr).Because(accessErr)
}

func main() {
	// 3. Call the function and get the rich, chained error.
	err := loadConfig("/etc/app/config.yaml")

	// 4. Print the full error chain for detailed logging. 🪵
	// The output is a clear, human-readable trace of what happened.
	fmt.Printf("Full error: %s\n\n", err)

	// 5. Check for specific errors to make programmatic decisions.
	// This works even though the errors are deeply nested in the chain.
	fmt.Println("Checking error identities...")

	if errors.Is(err, ErrConfigValidation) {
		fmt.Println("✅ High-level operation failed: Could not validate config.")
	}

	if errors.Is(err, ErrFileAccess) {
		fmt.Println("✅ Intermediate cause: Could not access the file.")
	}

	// You can even check against standard library errors wrapped in the chain!
	if errors.Is(err, os.ErrPermission) {
		fmt.Println("✅ Root cause: Permission was denied by the OS.")
	}
}

// Output:
// Full error: config validation failed: file access error: permission denied
//
// Checking error identities...
// ✅ High-level operation failed: Could not validate config.
// ✅ Intermediate cause: Could not access the file.
// ✅ Root cause: Permission was denied by the OS.
```

## Development

### System Requirements

```shell
go version
# go version go1.25.3

just --version
# just 1.42.4
```

### Download sources

```shell
PROJECT_ROOT=ex
git clone https://github.com/therenotomorrow/ex.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

### Setup dependencies

```shell
# install dependencies
go mod tidy

# check code integrity
just code test

# setup safe development (optional)
git config --local core.hooksPath .githooks
```

## Testing

```shell
# run quick checks
just test smoke # or just test

# run with coverage
just test cover
```

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.
