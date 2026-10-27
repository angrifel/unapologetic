package bytesx_test

import (
	"fmt"
	"io"

	"github.com/angrifel/unapologetic/bytesx"
)

func ExampleReadSeekBuffer() {
	var buf bytesx.ReadSeekBuffer
	buf.WriteString("Hello, World!")

	// Seek to an offset and do a partial read.
	buf.Seek(7, io.SeekStart)
	b := make([]byte, 5)
	buf.Read(b)
	fmt.Println(string(b))

	// Seek back to the start and read everything.
	buf.Seek(0, io.SeekStart)
	data, _ := io.ReadAll(&buf)
	fmt.Println(string(data))

	// Write new content; it appends to the buffer without resetting the read position.
	// Seek to position 7 and read from there, spanning both the original and new content.
	buf.WriteString(" Goodbye!")
	buf.Seek(7, io.SeekStart)
	data, _ = io.ReadAll(&buf)
	fmt.Println(string(data))
	// Output:
	// World
	// Hello, World!
	// World! Goodbye!
}

func ExampleReadSeekBuffer_ReadAt() {
	buf := bytesx.NewReadSeekBufferString("Hello, World!")

	b := make([]byte, 5)
	buf.ReadAt(b, 7)
	fmt.Println(string(b))

	// ReadAt does not advance the read position.
	data, _ := io.ReadAll(buf)
	fmt.Println(string(data))
	// Output:
	// World
	// Hello, World!
}

func ExampleReadSeekBuffer_Position() {
	buf := bytesx.NewReadSeekBufferString("Hello, World!")

	b := make([]byte, 5)
	buf.Read(b)
	fmt.Println(buf.Position())

	buf.Seek(7, io.SeekStart)
	fmt.Println(buf.Position())
	// Output:
	// 5
	// 7
}
