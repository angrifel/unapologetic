// Package errorsaux provides helpers for attaching stack traces to errors.
//
// Errors created with the standard library carry a message but no information about where they
// occurred. Without that information, the only way to know where a failure happened is to log it
// at every point it is returned, which is impractical and noisy. Attaching the stack trace to the
// error removes that need: the error can be returned up the call stack, wrapped with additional
// context along the way, and logged once wherever handling it is appropriate, while still
// reporting where the original failure occurred.
//
// # Capturing Stack Traces
//
// TraceError wraps an error with the stack trace of its caller:
//
//	if err := db.QueryRow(q).Scan(&v); err != nil {
//		return errorsaux.TraceError(err)
//	}
//
// TraceResult does the same for functions returning a value and an error, so their results can
// be returned directly:
//
//	return errorsaux.TraceResult(os.Open(name))
//
// Both return nil when there is no error. When the error already carries a stack trace anywhere
// in its chain, for example because it was traced and then wrapped with fmt.Errorf and %w, it is
// returned unchanged so that the stack trace of the original failure is preserved. Any error with
// a StackTrace() []uintptr method counts as carrying a stack trace, so errors from other packages
// that record their own stack traces are recognized too. Traced errors keep their message and
// support errors.Is, errors.As and errors.Unwrap.
//
// # Sentinel Errors
//
// Do not trace package-level sentinel errors such as var ErrNotFound = errors.New("not found").
// The stack trace would point to package initialization and, because existing stack traces are
// preserved, tracing the sentinel at the point where it is returned would have no effect. Return
// the sentinel traced at that point instead:
//
//	return errorsaux.TraceError(ErrNotFound)
//
// # Reading Stack Traces
//
// StackTrace returns the program counters recorded anywhere in an error's chain, or nil if there
// are none. FormatStackTrace renders them using the layout of Go panics, which is how the stack
// trace of a wrapped error is logged:
//
//	slog.Error("request failed", "err", err,
//		"stack", errorsaux.FormatStackTrace(errorsaux.StackTrace(err)))
//
// To process them in some other way, resolve them into functions, files and lines with
// runtime.CallersFrames:
//
//	frames := runtime.CallersFrames(errorsaux.StackTrace(err))
//	for {
//		frame, more := frames.Next()
//		fmt.Println(frame.Function, frame.File, frame.Line)
//		if !more {
//			break
//		}
//	}
//
// TracedError returns the error that carries the stack trace StackTrace reports. Joined errors,
// for example those created with errors.Join, can contain several traced errors, each with its own
// stack trace; TracedErrors returns all of them so each can be reported with its stack trace:
//
//	for _, e := range errorsaux.TracedErrors(err) {
//		slog.Error("request failed", "err", e,
//			"stack", errorsaux.FormatStackTrace(errorsaux.StackTrace(e)))
//	}
//
// Trace each error before joining it: TraceError returns an error unchanged when any error in it
// already carries a stack trace, so tracing the joined error does not record the others.
//
// Formatting a traced error with %+v prints its message followed by the stack trace, using the
// layout of Go panics and omitting frames from the Go runtime. All other verbs, flags and widths
// format the error message as if it were a string. Only the error being formatted controls its
// output: once a traced error is wrapped, for example with fmt.Errorf and %w, %+v prints the
// message alone, so use FormatStackTrace and StackTrace to print the stack trace.
package errorsaux
