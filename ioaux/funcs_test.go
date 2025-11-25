package ioaux

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestReaderFunc(t *testing.T) {
	// prepare test cases
	type testCase struct {
		rf        func(p []byte) (int, error)
		p         []byte
		bytesRead int
		err       error
	}

	rfErr := errors.New("test error")
	var testcases = []testCase{
		{
			rf:        func(p []byte) (int, error) { return 0, nil },
			p:         []byte{},
			bytesRead: 0,
			err:       nil,
		},
		{
			rf:        func(p []byte) (int, error) { return 0, io.EOF },
			p:         make([]byte, 5),
			bytesRead: 0,
			err:       io.EOF,
		},
		{
			rf:        func(p []byte) (int, error) { return copy(p, "Hello"), nil },
			p:         make([]byte, 10),
			bytesRead: 5,
			err:       nil,
		},
		{
			rf:        func(p []byte) (int, error) { return copy(p, "Hello"), nil },
			p:         make([]byte, 5),
			bytesRead: 5,
			err:       nil,
		},
		{
			rf:        func(p []byte) (int, error) { return 0, rfErr },
			p:         make([]byte, 10),
			bytesRead: 0,
			err:       rfErr,
		},
	}

	for tci, tc := range testcases {
		t.Run(strconv.Itoa(tci), func(t *testing.T) {
			// arrange
			rf := ReaderFunc(tc.rf)
			pSnapshot := make([]byte, len(tc.p))
			copy(pSnapshot, tc.p)

			// act
			br, err := rf.Read(tc.p)

			// assert
			assert.Equal(t, tc.bytesRead, br)
			assert.Equal(t, tc.err, err)
			if !bytes.Equal(pSnapshot[br:], tc.p[br:]) {
				t.Errorf("expected buffer tail %v, got %v", pSnapshot[br:], tc.p[br:])
			}
		})
	}
}

func TestCloserFunc(t *testing.T) {
	// prepare test cases
	type testCase struct {
		cf  func() error
		err error
	}

	cfErr := errors.New("close error")
	var testCases = []testCase{
		{
			cf:  func() error { return nil },
			err: nil,
		},
		{
			cf:  func() error { return cfErr },
			err: cfErr,
		},
	}

	for tci, tc := range testCases {
		t.Run(strconv.Itoa(tci), func(t *testing.T) {
			// arrange
			cf := CloserFunc(tc.cf)

			// act
			err := cf.Close()

			// assert
			assert.Equal(t, tc.err, err)
		})
	}
}
