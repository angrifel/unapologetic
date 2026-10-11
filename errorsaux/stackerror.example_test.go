package errorsaux_test

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/angrifel/unapologetic/errorsaux"
)

func ExampleTraceError() {
	_, err := os.Open("does-not-exist.txt")
	err = errorsaux.TraceError(err)

	fmt.Println(err)
	fmt.Println(errors.Is(err, os.ErrNotExist))
	// Output:
	// open does-not-exist.txt: no such file or directory
	// true
}

func ExampleTraceResult() {
	f, err := errorsaux.TraceResult(os.Open("does-not-exist.txt"))

	fmt.Println(f == nil)
	fmt.Println(err)
	fmt.Println(errorsaux.StackTrace(err) != nil)
	// Output:
	// true
	// open does-not-exist.txt: no such file or directory
	// true
}

func ExampleStackTrace() {
	err := fmt.Errorf("loading config: %w", errorsaux.TraceError(errors.New("boom")))

	frame, _ := runtime.CallersFrames(errorsaux.StackTrace(err)).Next()
	fmt.Println(frame.Function)
	// Output:
	// github.com/angrifel/unapologetic/errorsaux_test.ExampleStackTrace
}

func ExampleFormatStackTrace() {
	err := fmt.Errorf("loading config: %w", errorsaux.TraceError(errors.New("boom")))

	stack := errorsaux.FormatStackTrace(errorsaux.StackTrace(err))
	fmt.Println(strings.Split(stack, "\n")[1])
	// Output:
	// github.com/angrifel/unapologetic/errorsaux_test.ExampleFormatStackTrace
}

func ExampleTracedError() {
	traced := errorsaux.TraceError(errors.New("boom"))
	err := fmt.Errorf("loading config: %w", traced)

	fmt.Println(errorsaux.TracedError(err))
	fmt.Println(errorsaux.TracedError(err) == traced)
	// Output:
	// boom
	// true
}

func ExampleTracedErrors() {
	errWrite := errorsaux.TraceError(errors.New("write failed"))
	errClose := errorsaux.TraceError(errors.New("close failed"))
	err := fmt.Errorf("saving: %w", errors.Join(errWrite, errClose))

	for _, e := range errorsaux.TracedErrors(err) {
		frame, _ := runtime.CallersFrames(errorsaux.StackTrace(e)).Next()
		fmt.Println(e, "at", frame.Function)
	}
	// Output:
	// write failed at github.com/angrifel/unapologetic/errorsaux_test.ExampleTracedErrors
	// close failed at github.com/angrifel/unapologetic/errorsaux_test.ExampleTracedErrors
}
