package ioaux

import "io"

type readSeekerWithNopCloser struct {
	io.ReadSeeker
}

func (r *readSeekerWithNopCloser) Close() error { return nil }

var _ io.ReadSeekCloser = (*readSeekerWithNopCloser)(nil)

// ReadSeekerWithNopCloser wraps an io.ReadSeeker with a no-op Close method, returning an io.ReadSeekCloser.
func ReadSeekerWithNopCloser(readSeeker io.ReadSeeker) io.ReadSeekCloser {
	return &readSeekerWithNopCloser{ReadSeeker: readSeeker}
}
