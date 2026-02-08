package ioaux

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
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
		require.IsNil(t, err)
		assert.Equal(t, content, string(data))

		// Test seeking to beginning
		pos, err := rsc.Seek(0, io.SeekStart)
		require.IsNil(t, err)
		assert.Equal(t, int64(0), pos)

		// Read again after seeking
		data, err = io.ReadAll(rsc)
		require.IsNil(t, err)
		assert.Equal(t, content, string(data))

		// Test seeking from current position
		_, err = rsc.Seek(-5, io.SeekCurrent)
		require.IsNil(t, err)

		// Read partial content after seek
		buf := make([]byte, 5)
		n, err := rsc.Read(buf)
		require.IsNil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "orld!", string(buf))

		// Test seeking from end
		_, err = rsc.Seek(-6, io.SeekEnd)
		require.IsNil(t, err)

		buf = make([]byte, 6)
		n, err = rsc.Read(buf)
		require.IsNil(t, err)
		assert.Equal(t, 6, n)
		assert.Equal(t, "World!", string(buf))

		// Test close
		err = rsc.Close()
		assert.IsNil(t, err)
	})

	t.Run("read error propagation", func(t *testing.T) {
		expectedErr := errors.New("read error")
		src := io.NopCloser(iospy.ReaderWithEOFError(strings.NewReader("Hello W"), expectedErr))

		rsc := ReadSeekCloser(src)
		buffer := make([]byte, 5)

		// first Read should not return error
		n, err := rsc.Read(buffer)
		assert.IsNil(t, err)
		assert.Equal(t, "Hello", string(buffer[:n]))
		assert.Equal(t, 5, n)

		// Second Read should populate some additional bytes
		n, err = rsc.Read(buffer)
		assert.IsNil(t, err)
		assert.Equal(t, " W", string(buffer[:n]))
		assert.Equal(t, 2, n)

		// Third Read should return an error
		n, err = rsc.Read(buffer)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, "", string(buffer[:n]))
		assert.Equal(t, 0, n)

		// Close should not return the read error
		err = rsc.Close()
		assert.IsNil(t, err)
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
		require.IsNil(t, err)
		assert.Equal(t, "test", string(data))

		// Close should return the error
		err = rsc.Close()
		assert.Equal(t, expectedErr, err)
	})

	t.Run("empty reader", func(t *testing.T) {
		src := io.NopCloser(strings.NewReader(""))
		rsc := ReadSeekCloser(src)

		// Reading empty content
		data, err := io.ReadAll(rsc)
		require.IsNil(t, err)
		assert.Equal(t, 0, len(data))

		// Seeking in empty content
		pos, err := rsc.Seek(0, io.SeekStart)
		require.IsNil(t, err)
		assert.Equal(t, int64(0), pos)
	})

	t.Run("large content", func(t *testing.T) {
		// Create 1MB of test R
		content := bytes.Repeat([]byte("abcd"), 128*1024)
		src := io.NopCloser(bytes.NewReader(content))

		rsc := ReadSeekCloser(src)

		// Read all content
		data, err := io.ReadAll(rsc)
		require.IsNil(t, err)
		assert.EqualFunc(t, content, data, bytes.Equal)

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
			require.IsNil(t, err)

			buf := make([]byte, 10)
			n, err := rsc.Read(buf)
			require.IsNil(t, err)
			assert.Equal(t, 10, n)
			assert.Equal(t, pos.expect, string(buf))
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
		assert.Equal(t, readErr, err)

		// Close should fail with close error
		err = rsc.Close()
		assert.Equal(t, closeErr, err)
	})
}
