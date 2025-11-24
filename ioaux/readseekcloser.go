package ioaux

import (
	"bytes"
	"io"
)

var _ io.ReadSeekCloser = (*readSeekCloser)(nil)

type readSeekCloser struct {
	bytes.Reader
	readFromErr error
	closeErr    error
}

func (b *readSeekCloser) Close() error { return b.closeErr }
func (b *readSeekCloser) Read(p []byte) (n int, err error) {
	if n, err = b.Reader.Read(p); err == io.EOF && b.readFromErr != nil { //nolint:errorlint // the intention is to compare for io.EOF
		err = b.readFromErr
	}

	return
}

// ReadSeekCloser wraps an io.ReadCloser into an io.ReadSeekCloser, buffering the content
// in memory for seeking capabilities. if any error is returned from readCloser Read or Close,
// the same error will be returned on Read or Close respectively.
func ReadSeekCloser(readCloser io.ReadCloser) io.ReadSeekCloser {
	buf := &bytes.Buffer{}

	_, readFromErr := buf.ReadFrom(readCloser)
	closeErr := readCloser.Close()

	return &readSeekCloser{Reader: *bytes.NewReader(buf.Bytes()), readFromErr: readFromErr, closeErr: closeErr}
}
