package iospy

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestReadWithEOFError(t *testing.T) {
	t.Run("panic on nil error", func(t *testing.T) {
		reader := strings.NewReader("test data")

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, but did not panic")
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
			if err != customEOFErr {
				t.Errorf("expected error %v, got %v", customEOFErr, err)
			}
			if !bytes.Equal(buf, []byte{0, 0, 0, 0, 0}) {
				t.Errorf("expected %v, got %v", []byte{0, 0, 0, 0, 0}, buf)
			}
			if n != 0 {
				t.Errorf("expected 0 bytes read, got %d", n)
			}
		})

		t.Run("specified error on eof", func(t *testing.T) {
			reader := strings.NewReader("hello world")
			customEOFErr := errors.New("unkown error")

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

			if err1 != nil {
				t.Errorf("expected nil error, got %v", err1)
			}
			if err2 != nil {
				t.Errorf("expected nil error, got %v", err2)
			}
			if err3 != nil {
				t.Errorf("expected nil error, got %v", err3)
			}
			if err4 != customEOFErr {
				t.Errorf("expected error %v, got %v", customEOFErr, err4)
			}
			if n1 != 5 {
				t.Errorf("expected 5 bytes read, got %d", n1)
			}
			if n2 != 5 {
				t.Errorf("expected 5 bytes read, got %d", n2)
			}
			if n3 != 1 {
				t.Errorf("expected 1 byte read, got %d", n3)
			}
			if n4 != 0 {
				t.Errorf("expected 0 bytes read, got %d", n4)
			}
			if !bytes.Equal(buf1, []byte("hello")) {
				t.Errorf("expected %v, got %v", []byte("hello"), buf1)
			}
			if !bytes.Equal(buf2, []byte(" worl")) {
				t.Errorf("expected %v, got %v", []byte(" worl"), buf2)
			}
			if !bytes.Equal(buf3, []byte("d\x00\x00\x00\x00")) {
				t.Errorf("expected %v, got %v", []byte("d\x00\x00\x00\x00"), buf3)
			}
			if !bytes.Equal(buf4, []byte{0, 0, 0, 0, 0}) {
				t.Errorf("expected %v, got %v", []byte{0, 0, 0, 0, 0}, buf4)
			}

		})

		t.Run("underlying reader non-EOF error propagation", func(t *testing.T) {
			underlyingErr := errors.New("read failure")
			reader := ReaderWithEOFError(strings.NewReader("test"), underlyingErr)
			customEOFErr := errors.New("custom EOF")

			wrappedReader := ReaderWithEOFError(reader, customEOFErr)

			// Read available data
			buf := make([]byte, 10)
			n, err := wrappedReader.Read(buf)
			if err != nil {
				t.Errorf("expected nil error, got %v", err)
			}
			if n != 4 {
				t.Errorf("expected 4 bytes read, got %d", n)
			}
			if string(buf[:n]) != "test" {
				t.Errorf("expected %q, got %q", "test", string(buf[:n]))
			}

			n, err = wrappedReader.Read(buf)
			if err != underlyingErr {
				t.Errorf("expected error %v, got %v", underlyingErr, err)
			}
			if n != 0 {
				t.Errorf("expected 0 bytes read, got %d", n)
			}
			if string(buf[:n]) != "" {
				t.Errorf("expected empty string, got %q", string(buf[:n]))
			}
		})

	})

}
