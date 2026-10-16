// Package iospy provides testing utilities for observing and controlling io.Reader, io.Writer, and io.Closer behavior.
//
// This package offers tools for instrumenting I/O operations during testing, allowing you to:
//   - Record and inspect Read() and Close() method calls (witness pattern)
//   - Replace EOF errors with custom errors for testing error handling
//   - Create limited readers that return specific errors when limits are reached
//   - Create limited writers that return specific errors when limits are reached
//
// # Witness Types
//
// WitnessReader and WitnessCloser wrap standard io interfaces to record all method calls
// while preserving the original behavior. This is useful for verifying that code correctly
// handles I/O operations:
//
//	witnessed := iospy.WitnessReader(reader)
//	// Use witnessed as normal io.Reader
//	witnessed.Read(buf)
//	// Later inspect what happened
//	calls := witnessed.(iospy.ReaderWitness).ObservedReadCalls()
//
// # Error Control
//
// ReaderWithEOFError allows replacing EOF with custom errors to test error handling paths:
//
//	customErr := errors.New("custom error")
//	reader := iospy.ReaderWithEOFError(originalReader, customErr)
//	// When originalReader returns EOF, reader returns customErr instead
//
// LimitReaderWithError provides similar functionality to io.LimitReader but lets you
// specify what error to return when the limit is reached:
//
//	reader := iospy.LimitReaderWithError(r, 100, customErr)
//	// After 100 bytes, returns customErr instead of EOF
//
// LimitWriterWithError does the same for writes: it writes to an underlying writer and
// returns the specified error once n bytes have been written:
//
//	writer := iospy.LimitWriterWithError(w, 100, customErr)
//	// Writes beyond 100 bytes are cut off and return customErr
package iospy
