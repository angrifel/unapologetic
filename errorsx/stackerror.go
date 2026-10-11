package errorsx

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"slices"
	"strings"
)

// maxStackDepth is the maximum number of program counters recorded in a stack trace.
const maxStackDepth = 64

// Compile-time guards.
//
// The build fails if
//  1. T does not implement error.
//  2. T does not implement Unwrap.
//  3. T does not implement stackTracer.
var (
	_ error                       = (*tracedError)(nil)
	_ interface{ Unwrap() error } = (*tracedError)(nil)
	_ stackTracer                 = (*tracedError)(nil)
)

// stackTracer is implemented by errors that carry a stack trace.
type stackTracer interface {
	error
	StackTrace() []uintptr
}

// tracedError is used through a pointer so that errors compare by identity. Comparing values would
// compare the wrapped errors, which panics when their dynamic type is not comparable.
type tracedError struct {
	inner error
	pcs   []uintptr
}

func (e *tracedError) Error() string {
	return e.inner.Error()
}

func (e *tracedError) Unwrap() error {
	return e.inner
}

// StackTrace returns a copy of the recorded program counters, so callers cannot modify the error.
func (e *tracedError) StackTrace() []uintptr {
	return slices.Clone(e.pcs)
}

// StackTrace gets the stack trace from the first error in err's chain that has a StackTrace() []uintptr method,
// or nil if there is none.
func StackTrace(err error) []uintptr {
	var st stackTracer
	if errors.As(err, &st) {
		return st.StackTrace()
	}

	return nil
}

// TracedError returns the first error in err's tree that has a StackTrace() []uintptr method, or nil if
// there is none. The tree is searched in the same order as errors.As, so StackTrace(err) returns the stack
// trace of the error returned by TracedError(err).
func TracedError(err error) error {
	var st stackTracer
	if errors.As(err, &st) {
		return st
	}

	return nil
}

// TracedErrors returns every error in err's tree that has a StackTrace() []uintptr method, or nil if there
// are none. Use it to report the stack traces of joined errors, for example those created with errors.Join,
// where each joined error may carry its own stack trace: StackTrace returns the stack trace of each returned
// error.
//
// The tree is searched in the same order as errors.As, so the first returned error is the one returned by
// TracedError. Errors wrapped by a returned error are not searched: as with TraceError, the outermost stack
// trace on each path is the one that counts.
func TracedErrors(err error) []error {
	return appendTracedErrors(nil, err)
}

// appendTracedErrors appends the errors of err's tree that have a StackTrace() []uintptr method to dst,
// visiting err the way errors.As does: err itself, then its As method, then the errors it wraps.
func appendTracedErrors(dst []error, err error) []error {
	if err == nil {
		return dst
	}

	if st, ok := err.(stackTracer); ok { //nolint:errorlint // the tree is walked one error at a time.
		return append(dst, st)
	}

	if x, ok := err.(interface{ As(target any) bool }); ok {
		var st stackTracer
		if x.As(&st) {
			return append(dst, st)
		}
	}

	switch x := err.(type) { //nolint:errorlint // the tree is walked one error at a time.
	case interface{ Unwrap() error }:
		return appendTracedErrors(dst, x.Unwrap())
	case interface{ Unwrap() []error }:
		for _, inner := range x.Unwrap() {
			dst = appendTracedErrors(dst, inner)
		}
	}

	return dst
}

// FormatStackTrace formats program counters, as returned by StackTrace, using the layout of Go panics and
// omitting frames from the Go runtime. Each frame is preceded by a newline, so the result can be appended
// directly to an error message.
func FormatStackTrace(pcs []uintptr) string {
	if len(pcs) == 0 {
		return ""
	}

	var b strings.Builder

	frames := runtime.CallersFrames(pcs)

	for {
		f, more := frames.Next()
		if !strings.HasPrefix(f.Function, "runtime.") {
			_, _ = fmt.Fprintf(&b, "\n%s\n\t%s:%d", f.Function, f.File, f.Line)
		}

		if !more {
			return b.String()
		}
	}
}

func (e *tracedError) Format(s fmt.State, verb rune) {
	if verb != 'v' || !s.Flag('+') {
		_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), e.Error())

		return
	}

	_, _ = fmt.Fprintf(s, "%+v", e.inner) // let a wrapped error print its own detail
	_, _ = io.WriteString(s, FormatStackTrace(e.pcs))
}

// TraceError returns an error that captures the stack trace at the point of invocation, up to 64 frames deep.
//
// The error is wrapped with the stack trace.
// Returns the same error if err == nil or if any error in its chain has a StackTrace() []uintptr method.
func TraceError(err error) error {
	return traceErrorInternal(err)
}

// TraceResult returns a result and an error that captures the stack trace at the point of invocation, up to 64
// frames deep.
//
// The error is wrapped with the stack trace.
// Returns the same value untouched.
// Returns the same error if err == nil or if any error in its chain has a StackTrace() []uintptr method.
func TraceResult[T any](r T, err error) (T, error) {
	return r, traceErrorInternal(err)
}

func traceErrorInternal(err error) error {
	if err == nil {
		return nil
	}

	var st stackTracer
	if errors.As(err, &st) {
		return err
	}

	return stackErrorFrom(err)
}

func stackErrorFrom(err error) error {
	// callersSkip is the number of callers to skip when capturing the stack trace.
	// 1 = this function.
	// 2 = traceErrorInternal.
	// 3 = TraceError or TraceResult.
	// 4 = the caller of TraceError or TraceResult.
	const callersSkip = 4

	var pcs [maxStackDepth]uintptr

	n := runtime.Callers(callersSkip, pcs[:])

	return &tracedError{inner: err, pcs: slices.Clone(pcs[:n])}
}
