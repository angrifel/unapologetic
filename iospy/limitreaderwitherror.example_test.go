package iospy_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/angrifel/unapologetic/iospy"
)

func ExampleLimitReaderWithError() {
	limitErr := errors.New("limit exceeded")
	reader := iospy.LimitReaderWithError(strings.NewReader("Hello, World!"), 5, limitErr)

	buf := make([]byte, 10)

	// First read gets up to 5 bytes
	n, err := reader.Read(buf)
	fmt.Printf("n=%d err=%v data=%q\n", n, err, string(buf[:n]))

	// Second read hits the limit
	n, err = reader.Read(buf)
	fmt.Printf("n=%d err=%v\n", n, err)
	// Output:
	// n=5 err=<nil> data="Hello"
	// n=0 err=limit exceeded
}
