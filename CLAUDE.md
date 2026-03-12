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
- **internal/assert** - Internal test assertion helpers (for non-critical assertions)
- **internal/assertfunctions** - Internal shared assertion function implementations used by assert and require
- **internal/ioaux** - Internal I/O adapters: `ReaderFunc`/`CloserFunc` function adapters (internal counterpart of public ioaux)
- **internal/require** - Internal test assertion helpers (for critical assertions)
- **internal/netaux** - Network test utilities: embedded DNS server for resolver testing
- **internal/testaux** - Internal test utilities

Design principles: composable utilities, preserved error semantics, testability focus.

### Style conventions:
- **Test naming**: Tests follow the pattern `TestFunctionName_ScenarioDescription` (e.g., `TestBuildDNSResponse_MatchingARecord`)
- **Commit messages**: Use lowercase imperative style (e.g., "add feature", "fix bug", "simplify logic")
- **Doc comments**: Public symbols use Go-standard doc comments (`// FunctionName does X`)

### Opportunistic bug and typo reporting:
When reading files as part of any task, always report bugs and typos found outside the scope of the current task. For each finding, present it as a follow-up prompt at the end of the response using the `AskUserQuestion` tool, with options for corrective actions (e.g., fix it now, fix it later, ignore it).

### General coding principles:
- **No external dependencies**: Self-explanatory. Only take dependencies from golang's stdlib
- **No external services**: Features MUST NOT rely on calls to external services like web APIs or databases.
- **Data-driven as opposed to object-oriented**: Prefer data structures and functions over object-oriented design. In particular, the following sub-principles should be observed.
  - **Functions over struct with methods**: If a problem can be solved by a function or by a struct with methods, prefer the former over the latter.
  - **Data over state and behavior**: If a problem can be modeled with independent data and function calls, this should be the preferred approach to having structs with state and behavior.
  - **Public over private**: If a function or method can be independently used outside its intended purpose, said functionality should be public and not private. Only functions and methods that are both -- _single-purpose_ and _live within a constrained workflow_ -- should be made private.  

### Documentation:
Documentation is crucial to enable effective adoption of this code library. Here are some guidelines for writing good documentation:

- **Clear and concise**: Use simple language and avoid jargon. Make sure your documentation is easy to understand for anyone reading it.
- **Contextual**: Provide enough context for the reader to understand the purpose and usage of the code. Explain how the code fits into the overall system and what it is responsible for.
- **Package level documentation**: Include a high-level overview of the package's purpose and functionality. 
  - This should be written in doc.go for each public package as well as the doc.go at the root of the repository. 
  - Internal packages (i.e package prefixed with `internal/`) are exempt from this requirement.
  - The root `doc.go` must list all public packages.
  - Contents of the `doc.go` should be consistent with its package code.
- **Examples**: Include examples of how to use the code to help readers understand its usage and functionality.
  - Examples should be captured in go code in Example{FunctionName} format.
  - Examples should be as minimal as possible to demonstrate the intended usage.
  - Examples must be written in files named `{filename}.example_test.go`.
- **Version control**: Keep documentation up-to-date with changes to the codebase. Update documentation whenever you make changes to the code or add new features.

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
- Test cases should prefer table-driven tests to individual test cases.
- The assertions to be used are located in the `internal/assert` or `internal/require` package. If needed, use the standard library's `testing` package for assertions not found in `internal/assert` or `internal/require`
- Test cases should use real dependencies testing as opposed to mocks. In particular, the following should be observed at all times.
  - Whenever a test requires testing the outcome of I/O operations, it should prefer real I/O over mocked I/O
  - Whenever a test requires testing the outcome of HTTP operations, it should prefer real HTTP calls over mocked HTTP calls
  - Whenever a test requires testing DNS resolution or its metrics, it should override the DNS resolver used by the Dialer involved in the operations. If the latter were not possible, the default resolver should be overridden for the duration of the test.
  - When real I/O cannot produce a specific error condition needed for testing, request instructions on how to proceed before introducing wrappers or adapters. When such deviation is approved, The rationale and the alternatives considered must be written in the test cases so that later examination by humans or agents can effectively understand why such a decision was made. Agents can exercise liberty to document this in code comments in a way that is both intelligible to humans and friendly to itself

### Changelog conventions

- Changes to `internal/` packages are never changelog-worthy — they are not importable outside the module and are invisible to consumers.

### Testing exclusions

The following files are exempt from having unit tests:
- internal/require/require.go
- internal/assert/assert.go