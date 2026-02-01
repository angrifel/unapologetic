package iospy

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestLimitWriterWithError(t *testing.T) {
	t.Run("panic on nil error", func(t *testing.T) {
		var buf bytes.Buffer

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, but did not panic")
			}
		}()

		LimitWriterWithError(&buf, 5, nil)
	})

	t.Run("write cases", func(t *testing.T) {
		type chunkData struct {
			data                 []byte
			expectedBytesWritten int
			expectedError        error
		}

		var cd = func(data []byte, expectedBytesWritten int, expectedError error) chunkData {
			return chunkData{
				data:                 data,
				expectedBytesWritten: expectedBytesWritten,
				expectedError:        expectedError,
			}
		}

		leErr := errors.New("limit exceeded")
		var testCases = []struct {
			limit           int64
			limitError      error
			chunkData       []chunkData
			expectedContent string
		}{
			{
				limit:      0,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 0, leErr),
				},
				expectedContent: "",
			},
			{
				limit:      5,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 5, nil),
					cd([]byte(", World!"), 0, leErr),
				},
				expectedContent: "Hello",
			},
			{
				limit:      10,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 5, nil),
					cd([]byte(", W"), 3, nil),
					cd([]byte("orld!"), 2, leErr),
				},
				expectedContent: "Hello, Wor",
			},
			{
				limit:      13,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 5, nil),
					cd([]byte(", W"), 3, nil),
					cd([]byte("orld!"), 5, nil),
					cd([]byte("!!!"), 0, leErr),
				},
				expectedContent: "Hello, World!",
			},
			{
				limit:      13,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 5, nil),
					cd([]byte(", W"), 3, nil),
					cd([]byte("orld!"), 5, nil),
					cd([]byte("!!!"), 0, leErr),
					cd([]byte("x"), 0, leErr),
					cd([]byte("0123456789"), 0, leErr),
				},
				expectedContent: "Hello, World!",
			},
			{
				limit:      20,
				limitError: leErr,
				chunkData: []chunkData{
					cd([]byte("Hello"), 5, nil),
					cd([]byte(", W"), 3, nil),
					cd([]byte("orld!"), 5, nil),
				},
				expectedContent: "Hello, World!",
			},
		}

		for tci, tc := range testCases {
			t.Run(strconv.Itoa(tci), func(t *testing.T) {
				var buf bytes.Buffer
				limitedWriter := LimitWriterWithError(&buf, tc.limit, tc.limitError)
				for _, tccd := range tc.chunkData {
					n, err := limitedWriter.Write(tccd.data)
					assert.Equal(t, tccd.expectedError, err)
					assert.Equal(t, tccd.expectedBytesWritten, n)
				}
				assert.Equal(t, tc.expectedContent, buf.String())
			})
		}
	})

	t.Run("underlying writer error before limit", func(t *testing.T) {
		pr, pw := io.Pipe()
		pr.Close() // close read end so writes fail

		underlyingErr := errors.New("limit exceeded")
		limitedWriter := LimitWriterWithError(pw, 100, underlyingErr)

		_, err := limitedWriter.Write([]byte("Hello"))
		if err == nil {
			t.Fatal("expected error from underlying writer, got nil")
		}
		if errors.Is(err, underlyingErr) {
			t.Fatal("expected underlying writer error, not limit error")
		}
		pw.Close()
	})

	t.Run("different error types", func(t *testing.T) {
		testCases := []struct {
			name string
			err  error
		}{
			{"custom error", errors.New("custom limit error")},
			{"wrapped error", errors.New("wrapped: limit exceeded")},
			{"io error", io.ErrUnexpectedEOF},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				limitedWriter := LimitWriterWithError(&buf, 2, tc.err)

				// Write within limit
				n, err := limitedWriter.Write([]byte("Hi"))
				assert.IsNil(t, err)
				assert.Equal(t, 2, n)

				// Write beyond limit
				n, err = limitedWriter.Write([]byte("!!"))
				assert.Equal(t, tc.err, err)
				assert.Equal(t, 0, n)

				assert.Equal(t, "Hi", buf.String())
			})
		}
	})
}
