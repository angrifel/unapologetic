// Package unapologetic provides a collection of utility packages for Go development.
//
// unapologetic is organized into focused packages that extend the Go standard library
// with practical helpers for HTTP, I/O operations, and testing.
//
// # Packages
//
// httpaux - HTTP utilities for working with responses and round trippers:
//   - Clone and buffer HTTP response bodies
//   - Create RoundTripper implementations from functions
//   - Preserve error semantics when manipulating response bodies
//
// ioaux - I/O auxiliary utilities and adapters:
//   - Function adapters (ReaderFunc, CloserFunc) for implementing io interfaces
//   - ReadSeekCloser for adding seek capabilities to io.ReadCloser
//   - Memory-backed I/O wrappers with error preservation
//
// iospy - Testing utilities for observing and controlling I/O behavior:
//   - Witness wrappers for recording Read() and Close() calls
//   - Custom EOF error replacement for testing error paths
//   - Limited readers with configurable error behavior
//
// # Design Philosophy
//
// unapologetic packages are designed to be:
//   - Composable: Small, focused utilities that work well together
//   - Reliable: Preserve error semantics and panic behavior
//   - Testable: Provide visibility into I/O operations during testing
//   - Practical: Solve common real-world problems with minimal overhead
package unapologetic
