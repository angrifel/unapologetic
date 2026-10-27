package httpx

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
	"github.com/angrifel/unapologetic/iospy"
	"github.com/angrifel/unapologetic/iox"
)

func TestBufferResponse(t *testing.T) {
	t.Run("call with non-nil *http.Response", func(t *testing.T) {

		t.Run("with no read error and no close error", func(t *testing.T) {
			// arrange
			recorder := httptest.NewRecorder()
			recorder.Header().Set("Content-Type", "text/plain")
			recorder.WriteHeader(200)
			_, err := recorder.Write([]byte("Hello, World!"))
			require.IsNil(t, err)

			originalResponse := recorder.Result()
			readerWitness := iospy.WitnessReader(originalResponse.Body)
			closerWitness := iospy.WitnessCloser(originalResponse.Body)
			originalResponse.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: readerWitness,
				Closer: closerWitness,
			}

			// act
			resp := BufferResponseBody(originalResponse)

			// assert
			assert.NotEqual(t, nil, resp)

			// act
			bodyContent, bodyErr := io.ReadAll(resp.Body)

			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			readCallAfterReadingContent := len(readCalls)
			assert.Equal(t, nil, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			assert.GreaterOrEqual(t, readCallAfterReadingContent, 1)
			assert.Equal(t, io.EOF, readCalls[len(readCalls)-1].ResultErr)
			assert.Equal(t, 0, len(closeCalls))

			closeErr := resp.Body.Close()

			readCalls = readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls = closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			assert.Equal(t, nil, closeErr)
			assert.Equal(t, len(readCalls), readCallAfterReadingContent)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, nil, closeCalls[0].ResultErr)
		})

		t.Run("with read error and no close error", func(t *testing.T) {
			// arrange
			recorder := httptest.NewRecorder()
			recorder.Header().Set("Content-Type", "text/plain")
			recorder.WriteHeader(200)
			_, err := recorder.Write([]byte("Hello, World!"))
			require.IsNil(t, err)

			originalResponse := recorder.Result()
			errAtEOF := errors.New("read error")
			orbWithErr := iospy.ReaderWithEOFError(originalResponse.Body, errAtEOF)
			readerWitness := iospy.WitnessReader(orbWithErr)
			closerWitness := iospy.WitnessCloser(originalResponse.Body)
			originalResponse.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: readerWitness,
				Closer: closerWitness,
			}

			// act
			resp := BufferResponseBody(originalResponse)

			// assert
			assert.NotEqual(t, nil, resp)

			// act
			bodyContent, bodyErr := io.ReadAll(resp.Body)

			// assert
			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()
			readCallAfterReadingContent := len(readCalls)

			assert.Equal(t, errAtEOF, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			assert.GreaterOrEqual(t, readCallAfterReadingContent, 1)
			assert.Equal(t, errAtEOF, readCalls[len(readCalls)-1].ResultErr)
			assert.Equal(t, 0, len(closeCalls))

			// act
			closeErr := resp.Body.Close()
			readCalls = readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls = closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			// assert
			assert.Equal(t, nil, closeErr)
			assert.Equal(t, len(readCalls), readCallAfterReadingContent)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, nil, closeCalls[0].ResultErr)
		})

		t.Run("with no read error but close error", func(t *testing.T) {
			// arrange
			recorder := httptest.NewRecorder()
			recorder.Header().Set("Content-Type", "text/plain")
			recorder.WriteHeader(200)
			_, writeErr := recorder.Write([]byte("Hello, World!"))
			require.IsNil(t, writeErr)

			originalResponse := recorder.Result()
			orb := originalResponse.Body
			cErr := errors.New("close error")
			errCloser := iox.CloserFunc(func() error {
				if err := orb.Close(); err != nil {
					return err
				}

				return cErr
			})
			readerWitness := iospy.WitnessReader(orb)
			closerWitness := iospy.WitnessCloser(errCloser)
			originalResponse.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: readerWitness,
				Closer: closerWitness,
			}

			// act
			resp := BufferResponseBody(originalResponse)

			// assert
			assert.IsNotNil(t, resp)

			// act
			bodyContent, bodyErr := io.ReadAll(resp.Body)

			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			readCallAfterReadingContent := len(readCalls)
			require.Greater(t, readCallAfterReadingContent, 1)
			assert.Equal(t, io.EOF, readCalls[len(readCalls)-1].ResultErr)
			assert.IsNil(t, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))

			// act
			closeErr := resp.Body.Close()

			// assert
			readCalls = readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls = closerWitness.(iospy.CloserWitness).ObservedCloseCalls()
			assert.Equal(t, cErr, closeErr)
			assert.Equal(t, len(readCalls), readCallAfterReadingContent)
			require.Equal(t, 1, len(closeCalls))
			assert.Equal(t, cErr, closeCalls[0].ResultErr)
		})

		t.Run("with both read and close errors", func(t *testing.T) {
			// arrange
			recorder := httptest.NewRecorder()
			recorder.Header().Set("Content-Type", "text/plain")
			recorder.WriteHeader(200)
			_, err := recorder.Write([]byte("Hello, World!"))
			require.IsNil(t, err)

			originalResponse := recorder.Result()
			orb := originalResponse.Body
			errAtEOF := errors.New("read error")
			orbWithErr := iospy.ReaderWithEOFError(orb, errAtEOF)
			cErr := errors.New("close error")
			errCloser := iox.CloserFunc(func() error {
				if err := orb.Close(); err != nil {
					return err
				}

				return cErr
			})
			readerWitness := iospy.WitnessReader(orbWithErr)
			closerWitness := iospy.WitnessCloser(errCloser)
			originalResponse.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: readerWitness,
				Closer: closerWitness,
			}

			// act
			resp := BufferResponseBody(originalResponse)

			// assert
			require.IsNotNil(t, resp)

			bodyContent, bodyErr := io.ReadAll(resp.Body)

			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			readCallAfterReadingContent := len(readCalls)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			assert.Equal(t, errAtEOF, bodyErr)
			assert.GreaterOrEqual(t, len(readCalls), 1)
			assert.Equal(t, errAtEOF, readCalls[len(readCalls)-1].ResultErr)

			// act
			closeErr := resp.Body.Close()

			// assert
			readCalls = readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls = closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			assert.Equal(t, cErr, closeErr)
			assert.Equal(t, len(readCalls), readCallAfterReadingContent)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, cErr, closeCalls[0].ResultErr)
		})
	})
}
