package errorsx_test

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/angrifel/unapologetic/errorsx"
)

func ExampleTraceError() {
	_, err := os.Open("does-not-exist.txt")
	err = errorsx.TraceError(err)

	fmt.Println(err)
	fmt.Println(errors.Is(err, os.ErrNotExist))
	// Output:
	// open does-not-exist.txt: no such file or directory
	// true
}

func ExampleTraceResult() {
	f, err := errorsx.TraceResult(os.Open("does-not-exist.txt"))

	fmt.Println(f == nil)
	fmt.Println(err)
	fmt.Println(errorsx.StackTrace(err) != nil)
	// Output:
	// true
	// open does-not-exist.txt: no such file or directory
	// true
}

func ExampleStackTrace() {
	err := fmt.Errorf("loading config: %w", errorsx.TraceError(errors.New("boom")))

	frame, _ := runtime.CallersFrames(errorsx.StackTrace(err)).Next()
	fmt.Println(frame.Function)
	// Output:
	// github.com/angrifel/unapologetic/errorsx_test.ExampleStackTrace
}

func ExampleFormatStackTrace() {
	err := fmt.Errorf("loading config: %w", errorsx.TraceError(errors.New("boom")))

	stack := errorsx.FormatStackTrace(errorsx.StackTrace(err))
	fmt.Println(strings.Split(stack, "\n")[1])
	// Output:
	// github.com/angrifel/unapologetic/errorsx_test.ExampleFormatStackTrace
}

func ExampleTracedError() {
	traced := errorsx.TraceError(errors.New("boom"))
	err := fmt.Errorf("loading config: %w", traced)

	fmt.Println(errorsx.TracedError(err))
	fmt.Println(errorsx.TracedError(err) == traced)
	// Output:
	// boom
	// true
}

func ExampleTracedErrors() {
	errWrite := errorsx.TraceError(errors.New("write failed"))
	errClose := errorsx.TraceError(errors.New("close failed"))
	err := fmt.Errorf("saving: %w", errors.Join(errWrite, errClose))

	for _, e := range errorsx.TracedErrors(err) {
		frame, _ := runtime.CallersFrames(errorsx.StackTrace(e)).Next()
		fmt.Println(e, "at", frame.Function)
	}
	// Output:
	// write failed at github.com/angrifel/unapologetic/errorsx_test.ExampleTracedErrors
	// close failed at github.com/angrifel/unapologetic/errorsx_test.ExampleTracedErrors
}
