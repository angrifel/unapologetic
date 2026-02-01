# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Run all tests with race detection and coverage
make test

# Run linter (requires Docker)
make lint

# Tidy go modules
make tidy

# Run GitHub actions locally (requires act and Docker)
make run-github-action
```

## Architecture

**unapologetic** is a Go utility library (module: `github.com/angrifel/unapologetic`) organized into focused packages that extend the standard library:

- **httpaux** - HTTP utilities: response cloning (`CloneHTTPResponseWithBody`), body buffering (`BufferResponseBody`), and `RoundTripperFunc` adapter for creating `http.RoundTripper` from functions
- **httpspy** - HTTP message capture for debugging/logging with configurable header censoring and optional body collection
- **ioaux** - I/O adapters: `ReaderFunc`/`CloserFunc` function adapters, `ReadSeekCloser` for adding seek to `io.ReadCloser`
- **iospy** - Testing utilities: witness wrappers (`WitnessReader`, `WitnessCloser`) for recording I/O calls, `ReaderWithEOFError` for custom EOF replacement, `LimitReaderWithError`
- **internal/assert** - Internal test assertion helpers
- **internal/testaux** - Internal test utilities

Design principles: composable utilities, preserved error semantics, testability focus.

### General coding principles:
- **No external dependencies**: Self-explanatory. Only take dependencies from golang's stdlib
- **No external services**: Features MUST NOT rely on calls to external services like web APIs or databases.
- **Data-driven as opposed to object-oriented**: Prefer data structures and functions over object-oriented design. In particular the following sub-principles should be observed.
  - **Functions over struct with methods**: If a problem can be solved by a function or by a struct with methods, prefer the former over the latter.
  - **Data over state and behavior**: If a problem can be modeled with independent data and function calls, this should be the preferred approach to having structs with state and behavior.
  - **Public over private**: If a function or method can be independently used outside the intended purpose, said functionality should be public and not private. Only functions and methods that are both -- _single-purpose_ and _live within a constrained workflow_ -- should be made private.  

### Internal dependencies:
Some internal dependencies MUST NEVER be added to the project. The following is a exhausitve list of depdencies that must NEVER exist, the dependencie will be listed in the form `A -> B`, where `A` is the consuming package and `B` is the consumed package.

**Forbidden dependencies:**
- `iospy -> ioaux`
- `httpspy -> httpaux`

## Testing Instructions

Testing should be performed using the `make test` command, which runs all tests with race detection and coverage. Other methods of running test should be avoided as they are not through enough to provide adequate results.

Directives for writing unit tests:
- NEVER take dependencies on external libraries.
- NEVER take dependencies on external services.
- The assertions to be used are located in the `internal/assert` package. If needed, use the standard library's `testing` package for assertions not found in `internal/assert`
- Test cases should use real dependencies testing as opposed to mocks. In particular, the following should be observed at all times.
  - Whenever a test requires testing the outcome of I/O operations, it should prefer real I/O over mocked I/O
  - Whenever a test requires testing the outcome of HTTP operations, it should prefer real HTTP calls over mocked HTTP calls
  - Whenever a test requires testing DNS resolution or its metrics, it should override the DNS resolver used by the Dialer involved in the operations. If the latter were not possible, the default resolver should be overridden for the duration of the test.

