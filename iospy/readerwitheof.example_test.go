package iospy_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/angrifel/unapologetic/iospy"
)

func ExampleReaderWithEOFError() {
	customErr := errors.New("custom EOF error")
	reader := iospy.ReaderWithEOFError(strings.NewReader("hi"), customErr)

	buf := make([]byte, 10)

	// First read gets the data
	n, err := reader.Read(buf)
	fmt.Printf("n=%d err=%v\n", n, err)

	// Second read gets the custom error instead of io.EOF
	n, err = reader.Read(buf)
	fmt.Printf("n=%d err=%v\n", n, err)
	// Output:
	// n=2 err=<nil>
	// n=0 err=custom EOF error
}
