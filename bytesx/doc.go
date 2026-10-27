// Package bytesx fills gaps in the standard library's bytes package.
//
// ReadSeekBuffer is a variable-sized read/write buffer that combines bytes.Buffer and
// bytes.Reader. It provides the full bytes.Buffer API plus Seek, ReadAt, Size, and Position
// from bytes.Reader. Unlike bytes.Buffer, it keeps the bytes already read, so Seek and ReadAt
// can return to them; only Reset and Truncate discard them:
//
//	buf := bytesx.NewReadSeekBufferString("Hello, World!")
//	data, _ := io.ReadAll(buf)
//	_, _ = buf.Seek(0, io.SeekStart)
//	// buf can now be read again from the start
//
// ReadSeekBuffer is not safe for concurrent use.
package bytesx
