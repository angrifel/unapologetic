package iox

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
			assert.EqualFunc(t, pSnapshot[br:], tc.p[br:], bytes.Equal)
		})
	}
}

func TestSeekerFunc(t *testing.T) {
	// prepare test cases
	type testCase struct {
		result int64
		resErr error
		offset int64
		whence int
	}

	sfErr := errors.New("seek error")
	var testCases = []testCase{
		{
			result: 0,
			resErr: nil,
			offset: 0,
			whence: io.SeekStart,
		},
		{
			result: 15,
			resErr: nil,
			offset: 5,
			whence: io.SeekCurrent,
		},
		{
			result: 7,
			resErr: nil,
			offset: -3,
			whence: io.SeekEnd,
		},
		{
			result: 0,
			resErr: sfErr,
			offset: -1,
			whence: io.SeekStart,
		},
	}

	for tci, tc := range testCases {
		t.Run(strconv.Itoa(tci), func(t *testing.T) {
			// arrange
			var gotOffset int64
			var gotWhence int
			sf := SeekerFunc(func(offset int64, whence int) (int64, error) {
				gotOffset, gotWhence = offset, whence
				return tc.result, tc.resErr
			})

			// act
			pos, err := sf.Seek(tc.offset, tc.whence)

			// assert
			assert.Equal(t, tc.offset, gotOffset)
			assert.Equal(t, tc.whence, gotWhence)
			assert.Equal(t, tc.result, pos)
			assert.Equal(t, tc.resErr, err)
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
