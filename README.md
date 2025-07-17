# `ex`

Shows another way for working with the ~~ex~~ (:dancer:) errors in Go.

<div>
  <a href="https://github.com/therenotomorrow/ex/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/therenotomorrow/ex?color=FBC02D" alt="GitHub releases">
  </a>
  <a href="https://go.dev/doc/go1.21" target="_blank">
    <img src="https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg" alt="Go 1.21">
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

	"github.com/therenotomorrow/ex"
)

func main() {
	// Create a deeply nested error
	stdErr := errors.New("standard error")
	labelErr := ex.LError("ex const error")

	err := ex.New("xxx error").Because(ex.From(stdErr).Because(labelErr.Reason("why not?")))

	// Extract the original, underlying error
	cause := ex.Cause(err)
	fmt.Println(cause)       // why not?
	fmt.Println(err.Error()) // xxx error
}
```

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.

## Development

### System Requirements

```shell
go version
# go version go1.24.3

just --version
# just 1.40.0
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
go mod download
go mod verify

# check code integrity
just code test smoke

# setup safe development (optional)
git config --local core.hooksPath .githooks
```

Taste it :heart:
