// Package unapologetic provides a collection of utility packages for Go development.
//
// unapologetic is organized into focused packages that extend the Go standard library
// with practical helpers for HTTP, I/O operations, and testing.
//
// # Packages
//
// errorsaux - Error utilities for recording where errors originate:
//   - Capture the stack trace where an error first enters your code
//   - Keep the original stack trace when errors are wrapped further
//   - Print stack traces with %+v using the layout of Go panics
//   - Format stack traces of wrapped errors
//   - Find every traced error inside joined errors
//
// httpaux - HTTP utilities for working with responses and round trippers:
//   - Clone and buffer HTTP response bodies
//   - Create RoundTripper implementations from functions
//   - Preserve error semantics when manipulating response bodies
//
// httpspy - HTTP message capture for debugging, logging, and testing:
//   - Capture HTTP requests and responses with configurable header censoring
//   - Trace HTTP round trips with connection timing metrics
//   - Non-destructive snapshots that preserve original message state
//
// ioaux - I/O auxiliary utilities and adapters:
//   - Function adapters (ReaderFunc, SeekerFunc, CloserFunc) for implementing io interfaces
//   - ReadSeekCloser for adding seek capabilities to io.ReadCloser, buffering content as it is read
//   - ReadSeekerWithNopCloser for turning an io.ReadSeeker into an io.ReadSeekCloser
//   - ReadSeekBuffer, a read/write buffer combining bytes.Buffer and bytes.Reader
//
// iospy - Testing utilities for observing and controlling I/O behavior:
//   - Witness wrappers for recording Read() and Close() calls
//   - Custom EOF error replacement for testing error paths
//   - Limited readers and writers with configurable error behavior
//
// # Design Philosophy
//
// unapologetic packages are designed to be:
//   - Composable: Small, focused utilities that work well together
//   - Reliable: Preserve error semantics and panic behavior
//   - Testable: Provide visibility into I/O operations during testing
//   - Practical: Solve common real-world problems with minimal overhead
package unapologetic
