package iospy

import "io"

type readerWithEOFError struct {
	R io.Reader
	E error
}

func (e *readerWithEOFError) Read(p []byte) (n int, err error) {
	n, err = e.R.Read(p)
	if err == io.EOF {
		err = e.E
	}

	return n, err
}

// ReaderWithEOFError wraps an io.Reader and replaces the io.EOF error with a custom error.
func ReaderWithEOFError(reader io.Reader, eofError error) io.Reader {
	if eofError == nil {
		panic("eofError must not be nil")
	}

	return &readerWithEOFError{R: reader, E: eofError}
}
