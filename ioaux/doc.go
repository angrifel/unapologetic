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
//   - SeekerFunc: A function adapter that implements io.Seeker.
//   - ReadSeekCloser: Wraps an io.ReadCloser into an io.ReadSeekCloser. Content is read
//     from the wrapped io.ReadCloser on demand and buffered in memory as it is read, so
//     bytes already read can be read again by seeking. Seek works within the bytes read so
//     far and returns ErrSeekPastBytesRead for positions past them. Read and Close errors from the
//     wrapped io.ReadCloser are returned unchanged.
//   - ReadSeekerWithNopCloser: Wraps an io.ReadSeeker with a no-op Close method,
//     returning an io.ReadSeekCloser.
//   - ReadSeekBuffer: A variable-sized read/write buffer that combines bytes.Buffer and
//     bytes.Reader. It provides the full bytes.Buffer API plus Seek, ReadAt, Size, and Position.
//     Unlike bytes.Buffer, bytes already read are kept so Seek and ReadAt can return to them.
package ioaux
