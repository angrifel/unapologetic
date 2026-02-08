package iospy

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/ioaux"
	"github.com/angrifel/unapologetic/internal/require"
)

func TestReaderWitness(t *testing.T) {
	t.Run("captures successful reads", func(t *testing.T) {
		// arrange
		sr := strings.NewReader("hello world")
		rw := WitnessReader(sr).(*readerWitness)

		// act
		buf1 := make([]byte, 5)
		n1, err1 := rw.Read(buf1)

		buf2 := make([]byte, 6)
		n2, err2 := rw.Read(buf2)

		// assert
		assert.Equal(t, 5, n1)
		assert.IsNil(t, err1)
		assert.EqualFunc(t, []byte("hello"), buf1, bytes.Equal)

		assert.Equal(t, 6, n2)
		assert.IsNil(t, err2)
		assert.EqualFunc(t, []byte(" world"), buf2[:n2], bytes.Equal)

		calls := rw.ObservedReadCalls()
		require.Equal(t, 2, len(calls))

		assert.Equal(t, 5, calls[0].ResultN)
		assert.IsNil(t, calls[0].ResultErr)
		assert.IsNil(t, calls[0].PanicVal)
		assert.Equal(t, 5, len(calls[0].P))

		assert.Equal(t, 6, calls[1].ResultN)
		assert.IsNil(t, calls[1].ResultErr)
		assert.IsNil(t, calls[1].PanicVal)
		assert.Equal(t, 6, len(calls[1].P))
	})

	t.Run("captures reads with errors", func(t *testing.T) {
		// arrange
		expectedErr := errors.New("read error")
		faultyReader := ioaux.ReaderFunc(func(p []byte) (n int, err error) {
			p[0] = 'a'
			p[1] = 'b'
			p[2] = 'c'

			return 3, expectedErr
		})

		rw := WitnessReader(faultyReader).(*readerWitness)

		// act
		buf := make([]byte, 10)
		n, err := rw.Read(buf)

		// assert
		assert.Equal(t, 3, n)
		assert.Equal(t, expectedErr, err)
		assert.EqualFunc(t, []byte("abc\x00\x00\x00\x00\x00\x00\x00"), buf, bytes.Equal)

		calls := rw.ObservedReadCalls()
		require.Equal(t, 1, len(calls))
		assert.Equal(t, 3, calls[0].ResultN)
		assert.Equal(t, expectedErr, calls[0].ResultErr)
		assert.IsNil(t, calls[0].PanicVal)
		assert.Equal(t, 10, len(calls[0].P))
	})

	t.Run("captures panics and re-panics", func(t *testing.T) {
		// arrange
		expectedPanicVal := "read panic"
		fakeReader := ioaux.ReaderFunc(func(p []byte) (n int, err error) {
			panic(expectedPanicVal)
		})
		rw := WitnessReader(fakeReader).(*readerWitness)

		// act
		buf := make([]byte, 10)
		var panicVal interface{}
		func() {
			defer func() {
				panicVal = recover()
			}()
			_, _ = rw.Read(buf)
		}()

		// assert
		if panicVal == nil {
			t.Fatal("expected panic, but did not panic")
		}
		assert.Equal[any](t, expectedPanicVal, panicVal)

		calls := rw.ObservedReadCalls()

		require.Equal(t, 1, len(calls))
		assert.Equal(t, 0, calls[0].ResultN)
		assert.IsNil(t, calls[0].ResultErr)
		assert.Equal[any](t, expectedPanicVal, calls[0].PanicVal)
		assert.Equal(t, 10, len(calls[0].P))
	})
}
