package errorsaux

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
)

func TestTraceError(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		// arrange
		// nothing

		// act
		err := TraceError(nil)

		// assert
		assert.IsNil(t, err)
		assert.IsNil(t, StackTrace(err))
	})

	t.Run("with non-nil common error", func(t *testing.T) {
		// arrange
		err := errors.New("test error")
		_, file, line, ok := runtime.Caller(0)
		require.Equal(t, true, ok)

		// act
		wsErr := TraceError(err)

		// assert
		assert.NotEqual(t, err, wsErr)
		assert.Equal(t, fmt.Sprintf("%+v", err), fmt.Sprintf("%v", wsErr))
		assert.Equal(t, err, errors.Unwrap(wsErr))
		errorWithStackTrace := fmt.Sprintf("%+v", wsErr)
		stackTrace := StackTrace(wsErr)
		frames := runtime.CallersFrames(stackTrace)
		f1, more1 := frames.Next()
		assert.Equal(t, "github.com/angrifel/unapologetic/errorsaux.TestTraceError.func2", f1.Function)
		assert.Equal(t, stackTrace[0]-1, f1.PC)
		assert.Equal(t, file, f1.File)
		assert.Equal(t, line+4, f1.Line)
		assert.Equal(t, true, more1)
		f2, more2 := frames.Next()
		assert.Equal(t, "testing.tRunner", f2.Function)
		assert.Equal(t, stackTrace[1]-1, f2.PC)
		assert.Equal(t, true, more2)
		f3, more3 := frames.Next()
		assert.Equal(t, "runtime.goexit", f3.Function)
		assert.Equal(t, stackTrace[2]-1, f3.PC)
		assert.Equal(t, false, more3)
		expectedErrorWithStackTrace := err.Error() + "\n" +
			f1.Function + "\n\t" + f1.File + ":" + strconv.Itoa(f1.Line) + "\n" +
			f2.Function + "\n\t" + f2.File + ":" + strconv.Itoa(f2.Line)
		assert.Equal(t, expectedErrorWithStackTrace, errorWithStackTrace)
	})

	t.Run("with non-nil stackError", func(t *testing.T) {
		// arrange
		tracedErr := TraceError(errors.New("test error"))
		err := fmt.Errorf("wrapped: %w", tracedErr)
		_, file, line, ok := runtime.Caller(0)
		require.Equal(t, true, ok)

		// act
		wsErr := TraceError(err)

		// assert
		assert.IsNotNil(t, wsErr)
		assert.Equal(t, err, wsErr)
		assert.Equal(t, fmt.Sprintf("%+v", err), fmt.Sprintf("%v", wsErr))
		assert.Equal(t, tracedErr, errors.Unwrap(wsErr))
		errorWithStackTrace := fmt.Sprintf("%+v", tracedErr)
		stackTrace := StackTrace(wsErr)
		frames := runtime.CallersFrames(stackTrace)
		f1, more1 := frames.Next()
		assert.Equal(t, "github.com/angrifel/unapologetic/errorsaux.TestTraceError.func3", f1.Function)
		assert.Equal(t, stackTrace[0]-1, f1.PC)
		assert.Equal(t, file, f1.File)
		assert.Equal(t, line-2, f1.Line)
		assert.Equal(t, true, more1)
		f2, more2 := frames.Next()
		assert.Equal(t, "testing.tRunner", f2.Function)
		assert.Equal(t, stackTrace[1]-1, f2.PC)
		assert.Equal(t, true, more2)
		f3, more3 := frames.Next()
		assert.Equal(t, "runtime.goexit", f3.Function)
		assert.Equal(t, stackTrace[2]-1, f3.PC)
		assert.Equal(t, false, more3)
		expectedErrorWithStackTrace := tracedErr.Error() + "\n" +
			f1.Function + "\n\t" + f1.File + ":" + strconv.Itoa(f1.Line) + "\n" +
			f2.Function + "\n\t" + f2.File + ":" + strconv.Itoa(f2.Line)
		assert.Equal(t, expectedErrorWithStackTrace, errorWithStackTrace)
	})
}

func TestTraceResult(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		// arrange
		v := 1

		// act
		nv, err := TraceResult(v, nil)

		// assert
		assert.Equal(t, v, nv)
		assert.IsNil(t, err)
		assert.IsNil(t, StackTrace(err))
	})

	t.Run("with non-nil common error", func(t *testing.T) {
		// arrange
		v := 1
		err := errors.New("test error")
		_, file, line, ok := runtime.Caller(0)
		require.Equal(t, true, ok)

		// act
		nv, wsErr := TraceResult(v, err)

		// assert
		assert.Equal(t, v, nv)
		assert.NotEqual(t, err, wsErr)
		assert.Equal(t, fmt.Sprintf("%+v", err), fmt.Sprintf("%v", wsErr))
		assert.Equal(t, err, errors.Unwrap(wsErr))
		errorWithStackTrace := fmt.Sprintf("%+v", wsErr)
		stackTrace := StackTrace(wsErr)
		frames := runtime.CallersFrames(stackTrace)
		f1, more1 := frames.Next()
		assert.Equal(t, "github.com/angrifel/unapologetic/errorsaux.TestTraceResult.func2", f1.Function)
		assert.Equal(t, stackTrace[0]-1, f1.PC)
		assert.Equal(t, file, f1.File)
		assert.Equal(t, line+4, f1.Line)
		assert.Equal(t, true, more1)
		f2, more2 := frames.Next()
		assert.Equal(t, "testing.tRunner", f2.Function)
		assert.Equal(t, stackTrace[1]-1, f2.PC)
		assert.Equal(t, true, more2)
		f3, more3 := frames.Next()
		assert.Equal(t, "runtime.goexit", f3.Function)
		assert.Equal(t, stackTrace[2]-1, f3.PC)
		assert.Equal(t, false, more3)
		expectedErrorWithStackTrace := err.Error() + "\n" +
			f1.Function + "\n\t" + f1.File + ":" + strconv.Itoa(f1.Line) + "\n" +
			f2.Function + "\n\t" + f2.File + ":" + strconv.Itoa(f2.Line)
		assert.Equal(t, expectedErrorWithStackTrace, errorWithStackTrace)
	})

	t.Run("with non-nil stackError", func(t *testing.T) {
		// arrange
		v := 1
		tracedErr := TraceError(errors.New("test error"))
		err := fmt.Errorf("wrapped: %w", tracedErr)
		_, file, line, ok := runtime.Caller(0)
		require.Equal(t, true, ok)

		// act
		nv, wsErr := TraceResult(v, err)

		// assert
		assert.Equal(t, v, nv)
		assert.IsNotNil(t, wsErr)
		assert.Equal(t, err, wsErr)
		assert.Equal(t, fmt.Sprintf("%+v", err), fmt.Sprintf("%v", wsErr))
		assert.Equal(t, tracedErr, errors.Unwrap(wsErr))
		errorWithStackTrace := fmt.Sprintf("%+v", tracedErr)
		stackTrace := StackTrace(wsErr)
		frames := runtime.CallersFrames(stackTrace)
		f1, more1 := frames.Next()
		assert.Equal(t, "github.com/angrifel/unapologetic/errorsaux.TestTraceResult.func3", f1.Function)
		assert.Equal(t, stackTrace[0]-1, f1.PC)
		assert.Equal(t, file, f1.File)
		assert.Equal(t, line-2, f1.Line)
		assert.Equal(t, true, more1)
		f2, more2 := frames.Next()
		assert.Equal(t, "testing.tRunner", f2.Function)
		assert.Equal(t, stackTrace[1]-1, f2.PC)
		assert.Equal(t, true, more2)
		f3, more3 := frames.Next()
		assert.Equal(t, "runtime.goexit", f3.Function)
		assert.Equal(t, stackTrace[2]-1, f3.PC)
		assert.Equal(t, false, more3)
		expectedErrorWithStackTrace := tracedErr.Error() + "\n" +
			f1.Function + "\n\t" + f1.File + ":" + strconv.Itoa(f1.Line) + "\n" +
			f2.Function + "\n\t" + f2.File + ":" + strconv.Itoa(f2.Line)
		assert.Equal(t, expectedErrorWithStackTrace, errorWithStackTrace)
	})
}

// foreignTracedError is an error from outside this package that records its own stack trace.
//
// It is not a mock: errorsaux recognizes stack traces through the StackTrace() []uintptr method
// rather than through its own type, and this type is the smallest real implementation of that
// contract. No error type in the standard library carries a stack trace, so a test-local type is
// the only way to exercise the contract without taking an external dependency.
type foreignTracedError struct {
	pcs []uintptr
}

func (e *foreignTracedError) Error() string { return "foreign" }

func (e *foreignTracedError) StackTrace() []uintptr { return e.pcs }

func TestFormatStackTrace_EmptyInput(t *testing.T) {
	testCases := []struct {
		name string
		pcs  []uintptr
	}{
		{name: "nil", pcs: nil},
		{name: "empty", pcs: []uintptr{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			stack := FormatStackTrace(tc.pcs)

			// assert
			assert.Equal(t, "", stack)
		})
	}
}

func TestFormatStackTrace_MatchesPlusVFormatting(t *testing.T) {
	// arrange
	err := TraceError(errors.New("test error"))

	// act
	stack := FormatStackTrace(StackTrace(err))

	// assert
	assert.Equal(t, fmt.Sprintf("%+v", err), "test error"+stack)
	assert.Equal(t, true, strings.HasPrefix(stack, "\ngithub.com/angrifel/unapologetic/errorsaux.TestFormatStackTrace_MatchesPlusVFormatting\n\t"))
	assert.Equal(t, false, strings.Contains(stack, "\nruntime."))
}

func TestStackTrace_WrappedTracedError(t *testing.T) {
	tracedErr := TraceError(errors.New("test error"))
	testCases := []struct {
		name string
		err  error
	}{
		{name: "not wrapped", err: tracedErr},
		{name: "wrapped once", err: fmt.Errorf("outer: %w", tracedErr)},
		{name: "wrapped twice", err: fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", tracedErr))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			stackTrace := StackTrace(tc.err)

			// assert
			assert.EqualFunc(t, StackTrace(tracedErr), stackTrace, slices.Equal)
			assert.Equal(t, true, strings.HasPrefix(FormatStackTrace(stackTrace), "\ngithub.com/angrifel/unapologetic/errorsaux.TestStackTrace_WrappedTracedError\n\t"))
		})
	}
}

func TestStackTrace_ForeignStackTracer(t *testing.T) {
	foreignErr := &foreignTracedError{pcs: []uintptr{1, 2, 3}}
	testCases := []struct {
		name string
		err  error
	}{
		{name: "not wrapped", err: foreignErr},
		{name: "wrapped", err: fmt.Errorf("outer: %w", foreignErr)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			stackTrace := StackTrace(tc.err)

			// assert
			assert.EqualFunc(t, foreignErr.pcs, stackTrace, slices.Equal)
		})
	}
}

func TestTraceError_ForeignStackTracer(t *testing.T) {
	foreignErr := &foreignTracedError{pcs: []uintptr{1, 2, 3}}
	testCases := []struct {
		name string
		err  error
	}{
		{name: "not wrapped", err: foreignErr},
		{name: "wrapped", err: fmt.Errorf("outer: %w", foreignErr)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			tracedErr := TraceError(tc.err)
			_, tracedResultErr := TraceResult(1, tc.err)

			// assert
			assert.Equal(t, tc.err, tracedErr)
			assert.Equal(t, tc.err, tracedResultErr)
			assert.EqualFunc(t, foreignErr.pcs, StackTrace(tracedErr), slices.Equal)
		})
	}
}

// sliceError is an error whose dynamic type is not comparable, so comparing two of them with == panics.
type sliceError []string

func (e sliceError) Error() string { return strings.Join(e, ", ") }

func TestTraceError_ComparesByIdentity(t *testing.T) {
	testCases := []struct {
		name  string
		inner func() error
	}{
		{name: "comparable error", inner: func() error { return errors.New("test error") }},
		{name: "non-comparable error", inner: func() error { return sliceError{"test error"} }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			inner := tc.inner()
			traced1 := TraceError(inner)
			traced2 := TraceError(inner)

			// act
			isSelf := errors.Is(traced1, traced1)
			isOther := errors.Is(traced1, traced2)

			// assert
			assert.Equal(t, true, isSelf)
			assert.Equal(t, false, isOther)
			assert.Equal(t, true, traced1 == traced1) // identity is what is under test
			assert.Equal(t, false, traced1 == traced2)
		})
	}
}

func TestStackTrace_RecordedFramesOnly(t *testing.T) {
	// act
	stackTrace := StackTrace(TraceError(errors.New("test error")))

	// assert
	assert.Equal(t, 3, len(stackTrace)) // this test, testing.tRunner and runtime.goexit
	assert.Equal(t, false, slices.Contains(stackTrace, 0))
}

func TestStackTrace_ReturnsCopy(t *testing.T) {
	// arrange
	err := TraceError(errors.New("test error"))
	stackTrace := StackTrace(err)
	expected := slices.Clone(stackTrace)

	// act
	stackTrace[0] = 0

	// assert
	assert.EqualFunc(t, expected, StackTrace(err), slices.Equal)
}

func traceNested(depth int) error {
	if depth == 0 {
		return TraceError(errors.New("test error"))
	}

	return traceNested(depth - 1)
}

func TestStackTrace_MaxStackDepth(t *testing.T) {
	// act
	stackTrace := StackTrace(traceNested(2 * maxStackDepth))

	// assert
	assert.Equal(t, maxStackDepth, len(stackTrace))
	f1, _ := runtime.CallersFrames(stackTrace).Next()
	assert.Equal(t, "github.com/angrifel/unapologetic/errorsaux.traceNested", f1.Function)
}

// asTracedError exposes a stack trace only through its As method, the way errors.As lets an error
// match a target it is not assignable to.
//
// Like foreignTracedError, it is not a mock but the smallest real implementation of a contract under
// test: TracedError and TracedErrors must find stack traces exposed through As, as errors.As does. No
// error type in the standard library has such an As method.
type asTracedError struct {
	traced error
}

func (e asTracedError) Error() string { return "as " + e.traced.Error() }

func (e asTracedError) As(target any) bool {
	st, ok := target.(*stackTracer)
	if ok {
		*st, ok = e.traced.(stackTracer) //nolint:errorlint // the test controls the dynamic type.
	}

	return ok
}

func TestTracedErrors(t *testing.T) {
	tracedA := TraceError(errors.New("a"))
	tracedB := TraceError(errors.New("b"))
	tracedC := TraceError(errors.New("c"))
	foreignErr := &foreignTracedError{pcs: []uintptr{1, 2, 3}}
	exposedViaAs := asTracedError{traced: tracedA}
	// Built directly, because TraceError never wraps an error that already carries a stack trace.
	tracedInsideTraced := &tracedError{inner: tracedA, pcs: []uintptr{4, 5, 6}}
	testCases := []struct {
		name     string
		err      error
		expected []error
	}{
		{name: "nil", err: nil, expected: nil},
		{name: "untraced", err: errors.New("plain"), expected: nil},
		{name: "traced", err: tracedA, expected: []error{tracedA}},
		{name: "wrapped traced", err: fmt.Errorf("outer: %w", tracedA), expected: []error{tracedA}},
		{name: "joined traced", err: errors.Join(tracedA, tracedB), expected: []error{tracedA, tracedB}},
		{name: "joined with untraced", err: errors.Join(errors.New("plain"), tracedB), expected: []error{tracedB}},
		{name: "joined untraced", err: errors.Join(errors.New("x"), errors.New("y")), expected: nil},
		{name: "multiple %w", err: fmt.Errorf("%w, %w", tracedA, tracedB), expected: []error{tracedA, tracedB}},
		{
			name:     "nested joins",
			err:      errors.Join(tracedA, fmt.Errorf("inner: %w", errors.Join(tracedB, tracedC))),
			expected: []error{tracedA, tracedB, tracedC},
		},
		{name: "foreign", err: errors.Join(foreignErr, tracedB), expected: []error{foreignErr, tracedB}},
		{name: "via As", err: errors.Join(exposedViaAs, tracedB), expected: []error{tracedA, tracedB}},
		{name: "As without trace", err: asTracedError{traced: errors.New("plain")}, expected: nil},
		{name: "traced inside traced", err: tracedInsideTraced, expected: []error{tracedInsideTraced}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			tracedErrors := TracedErrors(tc.err)
			tracedError := TracedError(tc.err)

			// assert
			assert.EqualFunc(t, tc.expected, tracedErrors, slices.Equal)

			if len(tc.expected) == 0 {
				assert.IsNil(t, tracedError)
				assert.IsNil(t, StackTrace(tc.err))
			} else {
				assert.Equal(t, tc.expected[0], tracedError)
				assert.EqualFunc(t, StackTrace(tc.expected[0]), StackTrace(tc.err), slices.Equal)
			}

			for _, e := range tracedErrors {
				assert.IsNotNil(t, StackTrace(e))
			}
		})
	}
}
