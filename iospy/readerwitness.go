package iospy

import "io"

// ReaderWitness is an interface for objects that can provide information about Read method calls.
// It allows inspection of Read operation results including byte counts, errors, and any panics that occurred.
type ReaderWitness interface {
	// ObservedReadCalls returns a slice of all observed Read method calls with their inputs and results.
	ObservedReadCalls() []ObservedReadCallArgs
}

var _ ReaderWitness = (*readerWitness)(nil)

type readerWitness struct {
	inner io.Reader
	calls []ObservedReadCallArgs
}

// ObservedReadCallArgs contains information about a single Read method call.
// It records both the input buffer and all execution results, including any panic that might have occurred.
type ObservedReadCallArgs struct {
	// P is the byte slice that was passed to Read.
	P []byte
	// ResultN is the number of bytes read, as returned by Read.
	ResultN int
	// ResultErr is the error returned by the Read method, if any.
	ResultErr error
	// PanicVal contains the value from any panic that occurred during Read.
	// It will be nil if no panic occurred.
	PanicVal any
}

// Read reads data into the provided byte slice and returns the number of bytes read along with any error encountered.
// The method captures each call, including input, results, and any panic, for later inspection.
// It re-panics if a panic occurs during the read operation.
func (r *readerWitness) Read(p []byte) (n int, err error) {
	defer func() {
		panicVal := recover()
		r.calls = append(r.calls, ObservedReadCallArgs{
			P:         p,
			ResultN:   n,
			ResultErr: err,
			PanicVal:  panicVal,
		})

		if panicVal != nil {
			panic(panicVal)
		}
	}()

	return r.inner.Read(p)
}

// ObservedReadCalls returns a slice of ObservedReadCallArgs, recording all calls made to the Read method, including input and results.
func (r *readerWitness) ObservedReadCalls() []ObservedReadCallArgs {
	return r.calls
}

// WitnessReader wraps an io.Reader with instrumentation that records all calls to Read().
// The returned object implements both io.Reader and ReaderWitness interfaces.
//
// The original Reader's behavior is preserved - all data, errors, and panics are propagated
// exactly as they would be from the underlying Reader, but each call is recorded
// and can be inspected via the ReaderWitness interface.
//
// This is particularly useful for testing to verify that a Reader was used correctly,
// to inspect what data was requested, and to monitor the results including any errors.
//
// Example:
//
//	file, _ := os.Open("filename.txt")
//	witnessed := WitnessReader(file)
//
//	// Use normally as an io.Reader
//	buf := make([]byte, 1024)
//	n, _ := witnessed.Read(buf)
//
//	// Then inspect call history in tests
//	calls := witnessed.(ReaderWitness).ObservedReadCalls()
func WitnessReader(reader io.Reader) io.Reader {
	return &readerWitness{inner: reader, calls: nil}
}
