package iospy

import (
	"io"
)

type limitedWriterWithError struct {
	W io.Writer // underlying writer
	N int64     // max bytes remaining
	E error
}

func (l *limitedWriterWithError) Write(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, l.E
	}

	if int64(len(p)) > l.N {
		n, err = l.W.Write(p[0:l.N])
		l.N -= int64(n)

		if err == nil {
			err = l.E
		}

		return
	}

	n, err = l.W.Write(p)
	l.N -= int64(n)

	return
}

// LimitWriterWithError returns a Writer that writes to w but stops with a specific error after n bytes.
// This function is similar in spirit to io.LimitReader but for writes, with the ability to control
// the error once the limit has been reached.
//
// Invoking LimitWriterWithError with a nil error will cause this function to panic.
func LimitWriterWithError(w io.Writer, n int64, err error) io.Writer {
	if err == nil {
		panic("err must not be nil")
	}

	return &limitedWriterWithError{W: w, N: n, E: err}
}
