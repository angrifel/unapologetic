// Package ioaux provides auxiliary utilities for working with I/O operations.
//
// This package extends the standard library's io package with adapters and helper
// functions that simplify common I/O patterns and provide additional capabilities.
//
// Key features:
//   - ReaderFunc: A function adapter that implements io.Reader, allowing functions
//     to be used directly as readers without defining custom types.
//   - CloserFunc: A function adapter that implements io.Closer, enabling functions
//     to serve as closers for resource cleanup.
//   - ReadSeekCloser: Wraps an io.ReadCloser into an io.ReadSeekCloser by buffering
//     content in memory, adding seeking capabilities while preserving error semantics.
package ioaux
