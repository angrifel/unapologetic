package iospy

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestReadWithEOFError(t *testing.T) {
	t.Run("panic on nil error", func(t *testing.T) {
		reader := strings.NewReader("test data")

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, but did not panic")
			}
		}()

		ReaderWithEOFError(reader, nil)
	})

	t.Run("read cases", func(t *testing.T) {

		t.Run("empty reader returns custom EOF error", func(t *testing.T) {
			reader := strings.NewReader("")
			customEOFErr := errors.New("empty reader error")

			wrappedReader := ReaderWithEOFError(reader, customEOFErr)

			// The first read should immediately return the custom EOF error
			buf := make([]byte, 5)
			n, err := wrappedReader.Read(buf)
			assert.Equal(t, customEOFErr, err)
			assert.EqualFunc(t, []byte{0, 0, 0, 0, 0}, buf, bytes.Equal)
			assert.Equal(t, 0, n)
		})

		t.Run("specified error on eof", func(t *testing.T) {
			reader := strings.NewReader("hello world")
			customEOFErr := errors.New("unknown error")

			wrappedReader := ReaderWithEOFError(reader, customEOFErr)

			buf1 := make([]byte, 5)
			buf2 := make([]byte, 5)
			buf3 := make([]byte, 5)
			buf4 := make([]byte, 5)

			// act
			n1, err1 := wrappedReader.Read(buf1)
			n2, err2 := wrappedReader.Read(buf2)
			n3, err3 := wrappedReader.Read(buf3)
			n4, err4 := wrappedReader.Read(buf4)

			assert.IsNil(t, err1)
			assert.IsNil(t, err2)
			assert.IsNil(t, err3)
			assert.Equal(t, customEOFErr, err4)
			assert.Equal(t, 5, n1)
			assert.Equal(t, 5, n2)
			assert.Equal(t, 1, n3)
			assert.Equal(t, 0, n4)
			assert.EqualFunc(t, []byte("hello"), buf1, bytes.Equal)
			assert.EqualFunc(t, []byte(" worl"), buf2, bytes.Equal)
			assert.EqualFunc(t, []byte("d\x00\x00\x00\x00"), buf3, bytes.Equal)
			assert.EqualFunc(t, []byte{0, 0, 0, 0, 0}, buf4, bytes.Equal)

		})

		t.Run("underlying reader non-EOF error propagation", func(t *testing.T) {
			underlyingErr := errors.New("read failure")
			reader := ReaderWithEOFError(strings.NewReader("test"), underlyingErr)
			customEOFErr := errors.New("custom EOF")

			wrappedReader := ReaderWithEOFError(reader, customEOFErr)

			// Read available data
			buf := make([]byte, 10)
			n, err := wrappedReader.Read(buf)
			assert.IsNil(t, err)
			assert.Equal(t, 4, n)
			assert.Equal(t, "test", string(buf[:n]))

			n, err = wrappedReader.Read(buf)
			assert.Equal(t, underlyingErr, err)
			assert.Equal(t, 0, n)
			assert.Equal(t, "", string(buf[:n]))
		})

	})

}
