package iox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
	"github.com/angrifel/unapologetic/iospy"
)

func TestReadSeekCloser(t *testing.T) {
	t.Run("does not read source readCloser before read or close calls", func(t *testing.T) {
		// arrange
		rc := io.NopCloser(strings.NewReader("0123456789abcdef"))

		readerWitness := iospy.WitnessReader(rc)
		closerWitness := iospy.WitnessCloser(rc)
		var sourceRC = struct {
			io.Reader
			io.Closer
		}{
			Reader: readerWitness,
			Closer: closerWitness,
		}

		// act
		rsc := ReadSeekCloser(sourceRC)

		// assert
		readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
		closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

		assert.IsNotNil(t, rsc)
		assert.Equal(t, 0, len(readCalls))
		assert.Equal(t, 0, len(closeCalls))
	})

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

func TestReadSeekCloser_SeekWithinBytesReadSoFar(t *testing.T) {
	// readN reads n bytes with io.ReadFull and reports what it got.
	readN := func(rsc io.ReadSeekCloser, n int) string {
		buf := make([]byte, n)
		m, err := io.ReadFull(rsc, buf)
		return fmt.Sprintf("read %q %v", buf[:m], err)
	}
	// seek seeks and reports the resulting position and error.
	seek := func(rsc io.ReadSeekCloser, offset int64, whence int) string {
		pos, err := rsc.Seek(offset, whence)
		return fmt.Sprintf("seek %d %v", pos, err)
	}

	errSeekPast := ErrSeekPastBytesRead.Error()

	tests := []struct {
		name            string
		act             func(rsc io.ReadSeekCloser) []string
		want            []string
		wantSourceReads int
	}{
		{
			name: "seeking past the bytes read so far fails and keeps the position",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{seek(rsc, 3, io.SeekStart), readN(rsc, 2), seek(rsc, 0, io.SeekStart), readN(rsc, 1), readN(rsc, 3), seek(rsc, 2, io.SeekStart), readN(rsc, 2), readN(rsc, 2), readN(rsc, 2)}
			},
			want:            []string{"seek 0 " + errSeekPast, `read "ab" <nil>`, "seek 0 <nil>", `read "a" <nil>`, `read "bcd" <nil>`, "seek 2 <nil>", `read "cd" <nil>`, `read "ef" <nil>`, `read "" EOF`},
			wantSourceReads: 4,
		},
		{
			name: "seeking to the end of the bytes read so far continues from the source",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{readN(rsc, 2), seek(rsc, 0, io.SeekStart), seek(rsc, 2, io.SeekStart), readN(rsc, 2)}
			},
			want:            []string{`read "ab" <nil>`, "seek 0 <nil>", "seek 2 <nil>", `read "cd" <nil>`},
			wantSourceReads: 2,
		},
		{
			name: "seeking past the bytes read so far with SeekCurrent fails and keeps the position",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{readN(rsc, 1), seek(rsc, 2, io.SeekCurrent), readN(rsc, 2), seek(rsc, 1, io.SeekStart), readN(rsc, 2)}
			},
			want:            []string{`read "a" <nil>`, "seek 0 " + errSeekPast, `read "bc" <nil>`, "seek 1 <nil>", `read "bc" <nil>`},
			wantSourceReads: 2,
		},
		{
			name: "SeekEnd is relative to the bytes read so far",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{seek(rsc, 0, io.SeekEnd), readN(rsc, 3), seek(rsc, -1, io.SeekEnd), readN(rsc, 1), seek(rsc, 1, io.SeekEnd), readN(rsc, 1)}
			},
			want:            []string{"seek 0 <nil>", `read "abc" <nil>`, "seek 2 <nil>", `read "c" <nil>`, "seek 0 " + errSeekPast, `read "d" <nil>`},
			wantSourceReads: 2,
		},
		{
			name: "seeking to a negative position fails and keeps the position",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{readN(rsc, 2), seek(rsc, -1, io.SeekStart), readN(rsc, 1)}
			},
			want:            []string{`read "ab" <nil>`, "seek 0 " + "iox.ReadSeekCloser.Seek: negative position", `read "c" <nil>`},
			wantSourceReads: 2,
		},
		{
			name: "seeking with an invalid whence fails and keeps the position",
			act: func(rsc io.ReadSeekCloser) []string {
				return []string{readN(rsc, 2), seek(rsc, 0, 3), readN(rsc, 1)}
			},
			want:            []string{`read "ab" <nil>`, "seek 0 " + "iox.ReadSeekCloser.Seek: invalid whence", `read "c" <nil>`},
			wantSourceReads: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			source := iospy.WitnessReader(strings.NewReader("abcdef"))
			rsc := ReadSeekCloser(io.NopCloser(source))

			// act
			got := tt.act(rsc)

			// assert
			assert.EqualFunc(t, tt.want, got, slices.Equal[[]string])
			assert.Equal(t, tt.wantSourceReads, len(source.(iospy.ReaderWitness).ObservedReadCalls()))
		})
	}
}
