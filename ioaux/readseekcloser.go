package ioaux

import (
	"errors"
	"io"
)

var _ io.ReadSeekCloser = (*readSeekCloser)(nil)

// ErrSeekPastBytesRead is returned by the Seek method of a [ReadSeekCloser] when the requested
// position is past the bytes read so far.
var ErrSeekPastBytesRead = errors.New("ioaux.ReadSeekCloser.Seek: position past the bytes read so far")

// ErrSeekInvalidWhence is returned by the Seek method of a [ReadSeekCloser] when an invalid whence value is provided.
var ErrSeekInvalidWhence = errors.New("ioaux.ReadSeekCloser.Seek: invalid whence")

// ErrSeekNegativePosition is returned by the Seek method of a [ReadSeekCloser] when a negative position is provided.
var ErrSeekNegativePosition = errors.New("ioaux.ReadSeekCloser.Seek: negative position")

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

	// The p buffer has some capacity left, we fill remaining bytes from the source
	// and perform accounting & seeking on the actual number of bytes read.
	bytesRead := n
	n, err = b.source.Read(p[bytesRead:])
	_, _ = b.buf.Write(p[bytesRead : bytesRead+n])
	_, _ = b.buf.Seek(int64(n), io.SeekCurrent)
	n += bytesRead

	return
}

// Seek fails for positions past the bytes read so far, so that Read always continues the source
// where the bytes read so far end, and no source byte is skipped. The new position is validated before
// the buffer is touched, so a failed Seek leaves the buffer exactly as it was.
func (b *readSeekCloser) Seek(offset int64, whence int) (int64, error) {
	size := b.buf.Size()

	var base int64

	switch whence {
	case io.SeekStart:
		// do nothing
	case io.SeekCurrent:
		base = b.buf.Position()
	case io.SeekEnd:
		base = size
	default:
		return 0, ErrSeekInvalidWhence
	}

	if newPos := base + offset; newPos < 0 {
		return 0, ErrSeekNegativePosition
	} else if newPos > size {
		return 0, ErrSeekPastBytesRead
	}

	return b.buf.Seek(offset, whence)
}

func (b *readSeekCloser) Close() error { return b.source.Close() }

// ReadSeekCloser wraps an io.ReadCloser into an io.ReadSeekCloser, buffering the content
// in memory for seeking capabilities. If any error is returned from readCloser Read or Close methods,
// the same error will be returned on Read or Close respectively.
//
// The content is read from readCloser on demand, as it is needed, and Seek works within the bytes
// read so far: [io.SeekEnd] is relative to the end of the bytes read so far, not to the end of the
// content of readCloser. Seeking past the bytes read so far returns [ErrSeekPastBytesRead] and leaves
// the position unchanged. Reading at the end of the bytes read so far continues from readCloser.
func ReadSeekCloser(readCloser io.ReadCloser) io.ReadSeekCloser {
	var rsb ReadSeekBuffer

	return &readSeekCloser{source: readCloser, buf: &rsb}
}
