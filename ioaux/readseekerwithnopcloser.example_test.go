package ioaux_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/angrifel/unapologetic/ioaux"
)

func ExampleReadSeekerWithNopCloser() {
	rsc := ioaux.ReadSeekerWithNopCloser(strings.NewReader("Hello, World!"))
	data, _ := io.ReadAll(rsc)
	fmt.Println(string(data))
	_, _ = rsc.Seek(7, io.SeekStart)
	data, _ = io.ReadAll(rsc)
	fmt.Println(string(data))
	_ = rsc.Close()
	// Output:
	// Hello, World!
	// World!
}
