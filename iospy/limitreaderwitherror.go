package iospy

import (
	"io"
)

type limitedReaderWithError struct {
	R io.Reader // underlying reader
	N int64     // max bytes remaining
	E error
}

func (l *limitedReaderWithError) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, l.E
	}

	n, err = l.R.Read(p[0:min(l.N, int64(len(p)))])
	l.N -= int64(n)

	return
}

// LimitReaderWithError returns a Reader that reads from r but stops with a specific error after n bytes.
// This function is similar in spirit io.LimitReader but with the ability to control the error once
// the limit has been reached.
//
// Invoking LimitReaderWithError with an error = io.EOF should be the same as using io.LimitReader.
// Invoking LimitReaderWithError with a nil error will cause this function to panic.
func LimitReaderWithError(r io.Reader, n int64, err error) io.Reader {
	if err == nil {
		panic("err must not be nil")
	}

	return &limitedReaderWithError{R: r, N: n, E: err}
}
