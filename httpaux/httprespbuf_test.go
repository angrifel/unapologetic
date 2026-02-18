package httpaux

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
	"github.com/angrifel/unapologetic/ioaux"
	"github.com/angrifel/unapologetic/iospy"
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
			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			assert.NotEqual(t, nil, resp)
			assert.Greater(t, len(readCalls), 1)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, io.EOF, readCalls[len(readCalls)-1].ResultErr)
			assert.Equal(t, nil, closeCalls[0].ResultErr)

			bodyContent, bodyErr := io.ReadAll(resp.Body)

			assert.Equal(t, nil, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			closeErr := resp.Body.Close()

			assert.Equal(t, nil, closeErr)
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
			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			assert.NotEqual(t, nil, resp)
			assert.Greater(t, len(readCalls), 1)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, errAtEOF, readCalls[len(readCalls)-1].ResultErr)
			assert.Equal(t, nil, closeCalls[0].ResultErr)

			bodyContent, bodyErr := io.ReadAll(resp.Body)
			assert.Equal(t, errAtEOF, bodyErr)

			assert.Equal(t, "Hello, World!", string(bodyContent))

			closeErr := resp.Body.Close()
			assert.Equal(t, nil, closeErr)
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
			errCloser := ioaux.CloserFunc(func() error {
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
			readCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			assert.NotEqual(t, nil, resp)
			assert.Greater(t, len(readCalls), 1)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, io.EOF, readCalls[len(readCalls)-1].ResultErr)
			assert.Equal(t, cErr, closeCalls[0].ResultErr)

			bodyContent, bodyErr := io.ReadAll(resp.Body)
			assert.Equal(t, nil, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			closeErr := resp.Body.Close()
			assert.Equal(t, cErr, closeErr)
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
			errCloser := ioaux.CloserFunc(func() error {
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
			readerCalls := readerWitness.(iospy.ReaderWitness).ObservedReadCalls()
			closeCalls := closerWitness.(iospy.CloserWitness).ObservedCloseCalls()

			require.IsNotNil(t, resp)
			assert.Greater(t, len(readerCalls), 1)
			assert.Equal(t, 1, len(closeCalls))
			assert.Equal(t, errAtEOF, readerCalls[len(readerCalls)-1].ResultErr)
			assert.Equal(t, cErr, closeCalls[0].ResultErr)

			bodyContent, bodyErr := io.ReadAll(resp.Body)
			assert.Equal(t, errAtEOF, bodyErr)
			assert.Equal(t, "Hello, World!", string(bodyContent))
			closeErr := resp.Body.Close()
			assert.Equal(t, cErr, closeErr)
		})
	})
}
