package ioaux

import (
	"io"
)

var _ io.ReadSeekCloser = (*readSeekCloser)(nil)

type readSeekCloser struct {
	source io.ReadCloser
	buf    *ReadSeekBuffer
}

func (b *readSeekCloser) Read(p []byte) (n int, err error) {
	if b.buf.Len() > 0 {
		n, _ = b.buf.Read(p)
	}

	if n == len(p) {
		// We've filled the p buffer with bytes read before, we return immediately.
		return
	}

	// We've the p buffer has some capacity left, we fill remaining bytes from the source
	// and perform accounting & seeking on the actual number of bytes read.
	bytesRead := n
	n, err = b.source.Read(p[bytesRead:])
	_, _ = b.buf.Write(p[bytesRead : bytesRead+n])
	_, _ = b.buf.Seek(int64(n), io.SeekCurrent)
	n += bytesRead

	return
}

func (b *readSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return b.buf.Seek(offset, whence)
}

func (b *readSeekCloser) Close() error { return b.source.Close() }

// ReadSeekCloser wraps an io.ReadCloser into an io.ReadSeekCloser, buffering the content
// in memory for seeking capabilities. If any error is returned from readCloser Read or Close methods,
// the same error will be returned on Read or Close respectively.
func ReadSeekCloser(readCloser io.ReadCloser) io.ReadSeekCloser {
	var rsb ReadSeekBuffer

	return &readSeekCloser{source: readCloser, buf: &rsb}
}
