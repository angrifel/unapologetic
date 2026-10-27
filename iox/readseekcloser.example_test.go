package iox_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/angrifel/unapologetic/iox"
)

func ExampleReadSeekCloser() {
	src := io.NopCloser(strings.NewReader("Hello, World!"))
	rsc := iox.ReadSeekCloser(src)

	// Read all content
	data, _ := io.ReadAll(rsc)
	fmt.Println(string(data))

	// Seek back to offset 1 and read again
	_, _ = rsc.Seek(1, io.SeekStart)
	data, _ = io.ReadAll(rsc)
	fmt.Println(string(data))
	// Output:
	// Hello, World!
	// ello, World!
}
