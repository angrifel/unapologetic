package ioaux

import "io"

type (

	// ReaderFunc is a function type that implements the io.Reader interface.
	ReaderFunc func([]byte) (int, error)

	// CloserFunc is a function type that implements the io.Closer interface.
	CloserFunc func() error
)

var (
	_ io.Reader = (ReaderFunc)(nil)
	_ io.Closer = (CloserFunc)(nil)
)

func (rf ReaderFunc) Read(p []byte) (n int, err error) { return rf(p) }

// Close calls the underlying function returns any error it produces.
func (cf CloserFunc) Close() error { return cf() }
