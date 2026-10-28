package iox

import "github.com/angrifel/unapologetic/internal/iox"

type (
	// ReaderFunc is a function type that implements the io.Reader interface.
	ReaderFunc = iox.ReaderFunc
	// SeekerFunc is a function type that implements the io.Seeker interface.
	SeekerFunc = iox.SeekerFunc
	// CloserFunc is a function type that implements the io.Closer interface.
	CloserFunc = iox.CloserFunc
)
