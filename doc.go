// Package unapologetic provides a collection of focused packages for Go development.
//
// The x packages fill gaps in the Go standard library for bytes, HTTP and I/O.
// The spy packages observe and control HTTP and I/O behavior for debugging and testing.
//
// # Packages
//
// bytesx - Fills gaps in the bytes package:
//   - ReadSeekBuffer, a read/write buffer combining bytes.Buffer and bytes.Reader
//   - Keeps bytes already read, so Seek and ReadAt can return to them
//
// httpx - Fills gaps in the net/http package for responses and round trippers:
//   - Clone and buffer HTTP response bodies
//   - Create RoundTripper implementations from functions
//   - Preserve error semantics when manipulating response bodies
//
// httpspy - HTTP message capture for debugging, logging, and testing:
//   - Capture HTTP requests and responses with configurable header censoring
//   - Trace HTTP round trips with connection timing metrics
//   - Non-destructive snapshots that preserve original message state
//
// iox - Fills gaps in the io package with adapters and wrappers:
//   - Function adapters (ReaderFunc, SeekerFunc, CloserFunc) for implementing io interfaces
//   - ReadSeekCloser for adding seek capabilities to io.ReadCloser, buffering content as it is read
//   - ReadSeekerWithNopCloser for turning an io.ReadSeeker into an io.ReadSeekCloser
//
// iospy - Testing utilities for observing and controlling I/O behavior:
//   - Witness wrappers for recording Read() and Close() calls
//   - Custom EOF error replacement for testing error paths
//   - Limited readers and writers with configurable error behavior
//
// # Design Philosophy
//
// unapologetic packages are designed to be:
//   - Composable: Small, focused packages that work well together
//   - Reliable: Preserve error semantics and panic behavior
//   - Testable: Provide visibility into I/O operations during testing
//   - Practical: Solve common real-world problems with minimal overhead
package unapologetic
