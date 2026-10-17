package ioaux_test

import (
	"fmt"
	"io"

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

func ExampleSeekerFunc() {
	seeker := ioaux.SeekerFunc(func(offset int64, whence int) (int64, error) {
		// your custom io.Seeker logic here
		// Example:
		return offset, nil
	})

	pos, _ := seeker.Seek(7, io.SeekStart)
	fmt.Println(pos)
	// Output:
	// 7
}

func ExampleCloserFunc() {
	closer := ioaux.CloserFunc(func() error {
		// your custom io.Closer logic here
		return nil
	})

	_ = closer.Close()
	// Output:
	//
}
