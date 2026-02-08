package iospy

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestLimitReaderWithError(t *testing.T) {
	t.Run("panic on nil error", func(t *testing.T) {
		reader := strings.NewReader("test R")

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, but did not panic")
			}
		}()

		LimitReaderWithError(reader, 5, nil)
	})

	t.Run("read cases", func(t *testing.T) {
		type chunkData struct {
			chunkSizes           int64
			expectedBytesRead    int
			expectedError        error
			expectedChunkContent []byte
		}

		var cd = func(chunkSize int64, expectedBytesRead int, expectedError error, expectedContent []byte) chunkData {
			return chunkData{
				chunkSizes:           chunkSize,
				expectedBytesRead:    expectedBytesRead,
				expectedError:        expectedError,
				expectedChunkContent: expectedContent,
			}
		}

		leErr := errors.New("limit exceeded")
		var testCases = []struct {
			content    io.Reader
			limit      int64
			limitError error
			chunkData  []chunkData
		}{
			{
				content:    strings.NewReader(""),
				limit:      0,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 0, leErr, []byte("\x00\x00\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader(""),
				limit:      1,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 0, io.EOF, []byte("\x00\x00\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      0,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 0, leErr, []byte("\x00\x00\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      10,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 2, nil, []byte("or\x00\x00\x00")),
					cd(3, 0, leErr, []byte("\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      13,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 5, nil, []byte("orld!")),
					cd(3, 0, leErr, []byte("\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      13,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 5, nil, []byte("orld!")),
					cd(3, 0, leErr, []byte("\x00\x00\x00")),
					cd(1, 0, leErr, []byte("\x00")),
					cd(10, 0, leErr, []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      14,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 5, nil, []byte("orld!")),
					cd(3, 0, io.EOF, []byte("\x00\x00\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      14,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 5, nil, []byte("orld!")),
					cd(3, 0, io.EOF, []byte("\x00\x00\x00")),
					cd(1, 0, io.EOF, []byte("\x00")),
				},
			},
			{
				content:    strings.NewReader("Hello, World!"),
				limit:      20,
				limitError: leErr,
				chunkData: []chunkData{
					cd(5, 5, nil, []byte("Hello")),
					cd(3, 3, nil, []byte(", W")),
					cd(5, 5, nil, []byte("orld!")),
					cd(3, 0, io.EOF, []byte("\x00\x00\x00")),
				},
			},
		}

		for tci, tc := range testCases {
			t.Run(strconv.Itoa(tci), func(t *testing.T) {

				limitedReader := LimitReaderWithError(tc.content, tc.limit, tc.limitError)
				for _, tccd := range tc.chunkData {
					buf := make([]byte, tccd.chunkSizes)
					n, err := limitedReader.Read(buf)
					assert.Equal(t, tccd.expectedError, err)
					assert.Equal(t, tccd.expectedBytesRead, n)
					assert.EqualFunc(t, tccd.expectedChunkContent, buf, bytes.Equal)
				}
			})
		}
	})

	t.Run("underlying reader error before limit", func(t *testing.T) {
		underlyingErr := errors.New("underlying reader error")
		reader := ReaderWithEOFError(strings.NewReader("Hi"), underlyingErr)
		limitErr := errors.New("limit exceeded")

		limitedReader := LimitReaderWithError(reader, 10, limitErr)

		// Read should succeed for available R
		buf := make([]byte, 5)
		n, err := limitedReader.Read(buf)
		assert.IsNil(t, err)
		assert.Equal(t, 2, n)
		assert.Equal(t, "Hi", string(buf[:n]))

		n, err = limitedReader.Read(buf)
		assert.Equal(t, underlyingErr, err)
		assert.Equal(t, 0, n)
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
				reader := strings.NewReader("test")
				limitedReader := LimitReaderWithError(reader, 2, tc.err)

				// Read within limit
				buf := make([]byte, 2)
				n, err := limitedReader.Read(buf)
				assert.IsNil(t, err)
				assert.Equal(t, 2, n)

				// Read beyond limit
				n, err = limitedReader.Read(buf)
				assert.Equal(t, tc.err, err)
				assert.Equal(t, 0, n)
			})
		}
	})
}
