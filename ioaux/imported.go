package ioaux

import "github.com/angrifel/unapologetic/internal/ioaux"

type (
	// ReaderFunc is a function type that implements the io.Reader interface.
	ReaderFunc = ioaux.ReaderFunc
	// CloserFunc is a function type that implements the io.Closer interface.
	CloserFunc = ioaux.CloserFunc
)
