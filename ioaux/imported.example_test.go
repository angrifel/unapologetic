package ioaux_test

import (
	"fmt"

	"github.com/angrifel/unapologetic/ioaux"
)

func ExampleReaderFunc() {
	reader := ioaux.ReaderFunc(func(p []byte) (int, error) {
		// your custom io.Reader logic here
		// Example:
		return copy(p, "hello"), nil
	})

	buf := make([]byte, 5)
	n, _ := reader.Read(buf)
	fmt.Println(string(buf[:n]))
	// Output:
	// hello
}

func ExampleCloserFunc() {
	closer := ioaux.CloserFunc(func() error {
		// you custom io.Closer logic here
		return nil
	})

	_ = closer.Close()
	// Output:
	//
}
