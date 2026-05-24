package ioaux

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
)

func TestReadSeekerWithNopCloser(t *testing.T) {
	const content = "Hello, World!"

	t.Run("read preserves behavior", func(t *testing.T) {
		ref := strings.NewReader(content)
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))

		wantData, wantErr := io.ReadAll(ref)
		gotData, gotErr := io.ReadAll(wrapper)

		assert.Equal(t, wantErr, gotErr)
		assert.EqualFunc(t, wantData, gotData, bytes.Equal)
	})

	t.Run("seek and re-read preserves behavior", func(t *testing.T) {
		ref := strings.NewReader(content)
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))

		// first read
		wantData, wantErr := io.ReadAll(ref)
		gotData, gotErr := io.ReadAll(wrapper)
		require.Equal(t, wantErr, gotErr)
		require.EqualFunc(t, wantData, gotData, bytes.Equal)

		// seek to start
		wantPos, wantSeekErr := ref.Seek(0, io.SeekStart)
		gotPos, gotSeekErr := wrapper.Seek(0, io.SeekStart)
		require.Equal(t, wantSeekErr, gotSeekErr)
		require.Equal(t, wantPos, gotPos)

		// re-read after seeking
		wantData, wantErr = io.ReadAll(ref)
		gotData, gotErr = io.ReadAll(wrapper)
		assert.Equal(t, wantErr, gotErr)
		assert.EqualFunc(t, wantData, gotData, bytes.Equal)
	})

	t.Run("partial read and seek preserves behavior", func(t *testing.T) {
		ref := strings.NewReader(content)
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))

		refBuf := make([]byte, 5)
		wrapBuf := make([]byte, 5)

		// partial read
		wantN, wantErr := ref.Read(refBuf)
		gotN, gotErr := wrapper.Read(wrapBuf)
		require.Equal(t, wantErr, gotErr)
		require.Equal(t, wantN, gotN)
		require.EqualFunc(t, refBuf[:wantN], wrapBuf[:gotN], bytes.Equal)

		// seek backwards from current position
		wantPos, wantSeekErr := ref.Seek(-3, io.SeekCurrent)
		gotPos, gotSeekErr := wrapper.Seek(-3, io.SeekCurrent)
		require.Equal(t, wantSeekErr, gotSeekErr)
		require.Equal(t, wantPos, gotPos)

		// partial read after seek
		wantN, wantErr = ref.Read(refBuf)
		gotN, gotErr = wrapper.Read(wrapBuf)
		assert.Equal(t, wantErr, gotErr)
		assert.Equal(t, wantN, gotN)
		assert.EqualFunc(t, refBuf[:wantN], wrapBuf[:gotN], bytes.Equal)
	})

	t.Run("seek from end preserves behavior", func(t *testing.T) {
		ref := strings.NewReader(content)
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))

		_, _ = io.ReadAll(ref)
		_, _ = io.ReadAll(wrapper)

		wantPos, wantSeekErr := ref.Seek(-6, io.SeekEnd)
		gotPos, gotSeekErr := wrapper.Seek(-6, io.SeekEnd)
		require.Equal(t, wantSeekErr, gotSeekErr)
		require.Equal(t, wantPos, gotPos)

		wantData, wantReadErr := io.ReadAll(ref)
		gotData, gotReadErr := io.ReadAll(wrapper)
		assert.Equal(t, wantReadErr, gotReadErr)
		assert.EqualFunc(t, wantData, gotData, bytes.Equal)
	})

	t.Run("close always returns nil", func(t *testing.T) {
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))
		assert.IsNil(t, wrapper.Close())
	})

	t.Run("close returns nil after read and seek", func(t *testing.T) {
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))
		_, _ = io.ReadAll(wrapper)
		_, _ = wrapper.Seek(0, io.SeekStart)
		assert.IsNil(t, wrapper.Close())
	})

	t.Run("close is idempotent", func(t *testing.T) {
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(content))
		assert.IsNil(t, wrapper.Close())
		assert.IsNil(t, wrapper.Close())
		assert.IsNil(t, wrapper.Close())
	})

	t.Run("empty reader preserves behavior", func(t *testing.T) {
		ref := strings.NewReader("")
		wrapper := ReadSeekerWithNopCloser(strings.NewReader(""))

		wantData, wantErr := io.ReadAll(ref)
		gotData, gotErr := io.ReadAll(wrapper)
		require.Equal(t, wantErr, gotErr)
		require.EqualFunc(t, wantData, gotData, bytes.Equal)

		wantPos, wantSeekErr := ref.Seek(0, io.SeekStart)
		gotPos, gotSeekErr := wrapper.Seek(0, io.SeekStart)
		require.Equal(t, wantSeekErr, gotSeekErr)
		assert.Equal(t, wantPos, gotPos)

		assert.IsNil(t, wrapper.Close())
	})

	t.Run("read error propagates unchanged", func(t *testing.T) {
		// strings.NewReader never returns errors on Read, so we use an anonymous
		// struct composing a ReaderFunc (which injects the error) with strings.NewReader's
		// Seek method. This is the only way to test error propagation for Read without
		// a full mock. Considered using iospy.ReaderWithEOFError, but it only implements
		// io.Reader — not io.ReadSeeker — so it cannot be used directly here.
		expectedErr := errors.New("read error")
		seeker := strings.NewReader("")
		rs := struct {
			io.Reader
			io.Seeker
		}{
			Reader: ReaderFunc(func(p []byte) (int, error) { return 0, expectedErr }),
			Seeker: seeker,
		}

		wrapper := ReadSeekerWithNopCloser(rs)

		_, err := wrapper.Read(make([]byte, 4))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("seek error propagates unchanged", func(t *testing.T) {
		// strings.NewReader never returns errors on Seek, so we use a SeekerFunc
		// to inject the error.
		expectedErr := errors.New("seek error")
		rs := struct {
			io.Reader
			io.Seeker
		}{
			Reader: strings.NewReader("data"),
			Seeker: SeekerFunc(func(offset int64, whence int) (int64, error) { return 0, expectedErr }),
		}

		wrapper := ReadSeekerWithNopCloser(rs)

		_, err := wrapper.Seek(0, io.SeekStart)
		assert.Equal(t, expectedErr, err)
	})
}
