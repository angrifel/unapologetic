// Package bytesx fills gaps in the standard library's bytes package.
//
// ReadSeekBuffer is a variable-sized read/write buffer that combines bytes.Buffer and
// bytes.Reader. It provides the full bytes.Buffer API plus Seek, ReadAt and Size from
// bytes.Reader, and Position, which returns the read position without changing the buffer.
// Unlike bytes.Buffer, it keeps the bytes already read, so Seek and ReadAt can return to them;
// only Reset (or, equivalently, Truncate(0)) discards them:
//
//	buf := bytesx.NewReadSeekBufferString("Hello, World!")
//	_, _ = io.ReadAll(buf)
//	_, _ = buf.Seek(0, io.SeekStart)
//	// buf can now be read again from the start
//
// ReadSeekBuffer is not safe for concurrent use.
package bytesx
