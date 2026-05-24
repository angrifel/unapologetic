package ioaux

import "io"

type (

	// ReaderFunc is a function type that implements the io.Reader interface.
	ReaderFunc func([]byte) (int, error)

	// SeekerFunc is a function type that implements the io.Seeker interface.
	SeekerFunc func(offset int64, whence int) (int64, error)

	// CloserFunc is a function type that implements the io.Closer interface.
	CloserFunc func() error
)

var (
	_ io.Reader = (ReaderFunc)(nil)
	_ io.Seeker = (SeekerFunc)(nil)
	_ io.Closer = (CloserFunc)(nil)
)

// Read calls the underlying function and returns its result.
func (rf ReaderFunc) Read(p []byte) (n int, err error) { return rf(p) }

// Seek calls the underlying function and returns its result.
func (f SeekerFunc) Seek(offset int64, whence int) (int64, error) { return f(offset, whence) }

// Close calls the underlying function returns any error it produces.
func (cf CloserFunc) Close() error { return cf() }
