package iospy

import "io"

// CloserWitness is an interface for objects that can provide information about Close method calls.
// It allows inspection of Close operation results including any errors or panics that occurred.
type CloserWitness interface {
	// ObservedCloseCalls returns a slice of all observed Close method calls with their results.
	ObservedCloseCalls() []ObservedCloseCallArgs
}

var _ CloserWitness = (*closerWitness)(nil)

type closerWitness struct {
	inner io.Closer
	calls []ObservedCloseCallArgs
}

// ObservedCloseCallArgs contains information about a single Close method call.
// It records both normal execution results and any panic that might have occurred.
type ObservedCloseCallArgs struct {
	// ResultErr is the error returned by the Close method, if any.
	ResultErr error
	// PanicVal contains the value from any panic that occurred during Close.
	// It will be nil if no panic occurred.
	PanicVal any
}

// Close wraps the inner Closer's Close method, records its result or any panic, and re-panics if a panic occurs.
func (c *closerWitness) Close() (err error) {
	defer func() {
		panicVal := recover()
		c.calls = append(c.calls, ObservedCloseCallArgs{
			ResultErr: err,
			PanicVal:  panicVal,
		})

		if panicVal != nil {
			panic(panicVal)
		}
	}()

	return c.inner.Close()
}

// ObservedCloseCalls returns a slice of ObservedCloseCallArgs containing details of all recorded Close method calls.
func (c *closerWitness) ObservedCloseCalls() []ObservedCloseCallArgs {
	return c.calls
}

// WitnessCloser wraps an io.Closer with instrumentation that records all calls to Close().
// The returned object implements both io.Closer and CloserWitness interfaces.
//
// The original Closer's behavior is preserved - all errors and panics are propagated
// exactly as they would be from the underlying Closer, but each call is recorded
// and can be inspected via the CloserWitness interface.
//
// Example:
//
//	file, _ := os.Open("filename.txt")
//	witnessed := WitnessCloser(file)
//
//	// Use normally as an io.Closer
//	witnessed.Close()
//
//	// Then inspect call history in tests
//	calls := witnessed.(CloserWitness).ObservedCloseCalls()
func WitnessCloser(closer io.Closer) io.Closer {
	return &closerWitness{inner: closer, calls: nil}
}
