package ioaux

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/iospy"
)

func TestReadSeekCloser(t *testing.T) {
	t.Run("successful read and seek operations", func(t *testing.T) {
		// Test R
		content := "Hello, World!"
		src := io.NopCloser(strings.NewReader(content))

		// Create ReadSeekCloser
		rsc := ReadSeekCloser(src)

		// Test reading
		data, err := io.ReadAll(rsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != content {
			t.Errorf("expected %q, got %q", content, string(data))
		}

		// Test seeking to beginning
		pos, err := rsc.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pos != 0 {
			t.Errorf("expected position 0, got %d", pos)
		}

		// Read again after seeking
		data, err = io.ReadAll(rsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != content {
			t.Errorf("expected %q, got %q", content, string(data))
		}

		// Test seeking from current position
		_, err = rsc.Seek(-5, io.SeekCurrent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Read partial content after seek
		buf := make([]byte, 5)
		n, err := rsc.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 5 {
			t.Errorf("expected 5 bytes read, got %d", n)
		}
		if string(buf) != "orld!" {
			t.Errorf("expected %q, got %q", "orld!", string(buf))
		}

		// Test seeking from end
		_, err = rsc.Seek(-6, io.SeekEnd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		buf = make([]byte, 6)
		n, err = rsc.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 6 {
			t.Errorf("expected 6 bytes read, got %d", n)
		}
		if string(buf) != "World!" {
			t.Errorf("expected %q, got %q", "World!", string(buf))
		}

		// Test close
		err = rsc.Close()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("read error propagation", func(t *testing.T) {
		expectedErr := errors.New("read error")
		src := io.NopCloser(iospy.ReaderWithEOFError(strings.NewReader("Hello W"), expectedErr))

		rsc := ReadSeekCloser(src)
		buffer := make([]byte, 5)

		// first Read should not return error
		n, err := rsc.Read(buffer)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(buffer[:n]) != "Hello" {
			t.Errorf("expected %q, got %q", "Hello", string(buffer[:n]))
		}
		if n != 5 {
			t.Errorf("expected 5 bytes read, got %d", n)
		}

		// Second Read should populate some additional bytes
		n, err = rsc.Read(buffer)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(buffer[:n]) != " W" {
			t.Errorf("expected %q, got %q", " W", string(buffer[:n]))
		}
		if n != 2 {
			t.Errorf("expected 2 bytes read, got %d", n)
		}

		// Third Read should return an error
		n, err = rsc.Read(buffer)
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if string(buffer[:n]) != "" {
			t.Errorf("expected empty string, got %q", string(buffer[:n]))
		}
		if n != 0 {
			t.Errorf("expected 0 bytes read, got %d", n)
		}

		// Close should not return the read error
		err = rsc.Close()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("close error propagation", func(t *testing.T) {
		expectedErr := errors.New("close error")
		src := struct {
			io.Reader
			io.Closer
		}{
			Reader: strings.NewReader("test"),
			Closer: CloserFunc(func() error { return expectedErr }),
		}

		rsc := ReadSeekCloser(src)

		// Reading should work
		data, err := io.ReadAll(rsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "test" {
			t.Errorf("expected %q, got %q", "test", string(data))
		}

		// Close should return the error
		err = rsc.Close()
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("empty reader", func(t *testing.T) {
		src := io.NopCloser(strings.NewReader(""))
		rsc := ReadSeekCloser(src)

		// Reading empty content
		data, err := io.ReadAll(rsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 0 {
			t.Errorf("expected empty slice, got %d bytes", len(data))
		}

		// Seeking in empty content
		pos, err := rsc.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pos != 0 {
			t.Errorf("expected position 0, got %d", pos)
		}
	})

	t.Run("large content", func(t *testing.T) {
		// Create 1MB of test R
		content := bytes.Repeat([]byte("abcd"), 128*1024)
		src := io.NopCloser(bytes.NewReader(content))

		rsc := ReadSeekCloser(src)

		// Read all content
		data, err := io.ReadAll(rsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(content, data) {
			t.Errorf("expected %d bytes, got %d bytes", len(content), len(data))
		}

		// Seek to random positions and verify content
		positions := []struct {
			offset int64
			whence int
			expect string
		}{
			{1024, io.SeekStart, string(content[1024:1034])},
			{-1024, io.SeekEnd, string(content[len(content)-1024 : len(content)-1014])},
			{100, io.SeekCurrent, string(content[len(content)-914 : len(content)-904])},
		}

		for _, pos := range positions {
			_, err := rsc.Seek(pos.offset, pos.whence)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			buf := make([]byte, 10)
			n, err := rsc.Read(buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != 10 {
				t.Errorf("expected 10 bytes read, got %d", n)
			}
			if string(buf) != pos.expect {
				t.Errorf("expected %q, got %q", pos.expect, string(buf))
			}
		}
	})

	t.Run("both read and close errors", func(t *testing.T) {
		readErr := errors.New("read error")
		closeErr := errors.New("close error")

		src := struct {
			io.Reader
			io.Closer
		}{
			Reader: iospy.ReaderWithEOFError(strings.NewReader(""), readErr),
			Closer: CloserFunc(func() error { return closeErr }),
		}

		rsc := ReadSeekCloser(src)

		// Read should fail with read error
		_, err := rsc.Read(make([]byte, 1))
		if !errors.Is(err, readErr) {
			t.Errorf("expected error %v, got %v", readErr, err)
		}

		// Close should fail with close error
		err = rsc.Close()
		if !errors.Is(err, closeErr) {
			t.Errorf("expected error %v, got %v", closeErr, err)
		}
	})
}
