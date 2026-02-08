package ioaux_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/angrifel/unapologetic/ioaux"
)

func ExampleReadSeekCloser() {
	src := io.NopCloser(strings.NewReader("Hello, World!"))
	rsc := ioaux.ReadSeekCloser(src)

	// Read all content
	data, _ := io.ReadAll(rsc)
	fmt.Println(string(data))

	// Seek back to start and read again
	_, _ = rsc.Seek(1, io.SeekStart)
	data, _ = io.ReadAll(rsc)
	fmt.Println(string(data))
	// Output:
	// Hello, World!
	// ello, World!
}
