set dotenv-load := false
set shell := ["sh", "-cu"]

BIN := justfile_directory() / ".bin"

[private]
default:
    @just --list --unsorted

# ---- golangci-lint

GOLANGCI_LINT_VERSION := 'v2.8.0'
GOLANGCI_LINT_PATH := BIN / 'golangci-lint'
GOLANGCI_LINT := GOLANGCI_LINT_PATH + '@' + GOLANGCI_LINT_VERSION

[private]
[doc('https://github.com/golangci/golangci-lint')]
install-golangci-lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b {{ BIN }} {{ GOLANGCI_LINT_VERSION }}
    mv {{ GOLANGCI_LINT_PATH }} {{ GOLANGCI_LINT }}

[doc('Run static analysis using `golangci-lint` to detect code issues')]
[group('code')]
lint:
    @if test ! -e {{ GOLANGCI_LINT }}; then just install-golangci-lint; fi
    {{ GOLANGCI_LINT }} run ./...

# ---- betteralign

BETTERALIGN_VERSION := 'v0.8.3'
BETTERALIGN_PATH := BIN / 'betteralign'
BETTERALIGN := BETTERALIGN_PATH + '@' + BETTERALIGN_VERSION

[private]
[doc('https://github.com/dkorunic/betteralign')]
install-betteralign:
    GOBIN={{ BIN }} go install github.com/dkorunic/betteralign/cmd/betteralign@{{ BETTERALIGN_VERSION }}
    mv {{ BETTERALIGN_PATH }} {{ BETTERALIGN }}

[doc('Reorder struct fields using `betteralign` to improve memory layout')]
[group('code')]
align:
    @if test ! -e {{ BETTERALIGN }}; then just install-betteralign; fi
    {{ BETTERALIGN }} -apply ./...

# ---- testing

[private]
smoke:
    go test ./...

[private]
cover:
    go test -count 1 -parallel 8 -race -coverprofile=coverage.out
    go tool cover -func coverage.out

# ---- shortcuts

[doc('Run all code quality tools')]
[group('code')]
code: align lint

[doc('Run tests by type: `smoke` for quick checks, `cover` for detailed analysis')]
[group('test')]
test type='smoke':
    just {{ type }}
