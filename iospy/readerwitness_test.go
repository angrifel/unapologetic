package iospy

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/ioaux"
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
		if n1 != 5 {
			t.Errorf("expected 5 bytes read, got %d", n1)
		}
		if err1 != nil {
			t.Errorf("expected nil error, got %v", err1)
		}
		if !bytes.Equal(buf1, []byte("hello")) {
			t.Errorf("expected %v, got %v", []byte("hello"), buf1)
		}

		if n2 != 6 {
			t.Errorf("expected 6 bytes read, got %d", n2)
		}
		if err2 != nil {
			t.Errorf("expected nil error, got %v", err2)
		}
		if !bytes.Equal(buf2[:n2], []byte(" world")) {
			t.Errorf("expected %v, got %v", []byte(" world"), buf2[:n2])
		}

		calls := rw.ObservedReadCalls()
		if len(calls) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(calls))
		}

		if calls[0].ResultN != 5 {
			t.Errorf("expected 5 bytes read, got %d", calls[0].ResultN)
		}
		if calls[0].ResultErr != nil {
			t.Errorf("expected nil error, got %v", calls[0].ResultErr)
		}
		if calls[0].PanicVal != nil {
			t.Errorf("expected nil panic, got %v", calls[0].PanicVal)
		}
		if len(calls[0].P) != 5 {
			t.Errorf("expected 5 bytes in buffer, got %d", len(calls[0].P))
		}

		if calls[1].ResultN != 6 {
			t.Errorf("expected 6 bytes read, got %d", calls[1].ResultN)
		}
		if calls[1].ResultErr != nil {
			t.Errorf("expected nil error, got %v", calls[1].ResultErr)
		}
		if calls[1].PanicVal != nil {
			t.Errorf("expected nil panic, got %v", calls[1].PanicVal)
		}
		if len(calls[1].P) != 6 {
			t.Errorf("expected 6 bytes in buffer, got %d", len(calls[1].P))
		}
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
		if n != 3 {
			t.Errorf("expected 3 bytes read, got %d", n)
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if !bytes.Equal(buf, []byte("abc\x00\x00\x00\x00\x00\x00\x00")) {
			t.Errorf("expected %v, got %v", []byte("abc\x00\x00\x00\x00\x00\x00\x00"), buf)
		}

		calls := rw.ObservedReadCalls()
		if len(calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(calls))
		}
		if calls[0].ResultN != 3 {
			t.Errorf("expected 3 bytes read, got %d", calls[0].ResultN)
		}
		if calls[0].ResultErr != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, calls[0].ResultErr)
		}
		if calls[0].PanicVal != nil {
			t.Errorf("expected nil panic, got %v", calls[0].PanicVal)
		}
		if len(calls[0].P) != 10 {
			t.Errorf("expected 10 bytes in buffer, got %d", len(calls[0].P))
		}
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
			t.Error("expected panic, but did not panic")
		}
		if panicVal != expectedPanicVal {
			t.Errorf("expected panic value %v, got %v", expectedPanicVal, panicVal)
		}

		calls := rw.ObservedReadCalls()

		if len(calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(calls))
		}
		if calls[0].ResultN != 0 {
			t.Errorf("expected 0 bytes read, got %d", calls[0].ResultN)
		}
		if calls[0].ResultErr != nil {
			t.Errorf("expected nil error, got %v", calls[0].ResultErr)
		}
		if calls[0].PanicVal != expectedPanicVal {
			t.Errorf("expected panic value %v, got %v", expectedPanicVal, calls[0].PanicVal)
		}
		if len(calls[0].P) != 10 {
			t.Errorf("expected 10 bytes in buffer, got %d", len(calls[0].P))
		}
	})
}
