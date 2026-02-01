package iospy_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/angrifel/unapologetic/iospy"
)

func ExampleLimitWriterWithError() {
	var buf bytes.Buffer
	limitErr := errors.New("limit exceeded")
	writer := iospy.LimitWriterWithError(&buf, 5, limitErr)

	// First write within limit
	n, err := writer.Write([]byte("Hello"))
	fmt.Printf("n=%d err=%v\n", n, err)

	// Second write exceeds limit
	n, err = writer.Write([]byte("!"))
	fmt.Printf("n=%d err=%v\n", n, err)

	fmt.Println(buf.String())
	// Output:
	// n=5 err=<nil>
	// n=0 err=limit exceeded
	// Hello
}
