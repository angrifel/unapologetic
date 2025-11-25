package httpaux

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/ioaux"
	"github.com/angrifel/unapologetic/iospy"
)

func TestBufferResponse(t *testing.T) {
	t.Run("call with non-nil error should echo the error back", func(t *testing.T) {
		resp := BufferResponseBody(nil)

		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("call with non-nil *http.Response", func(t *testing.T) {

		t.Run("with no read error and no close error", func(t *testing.T) {
			// arrange
			recorder := httptest.NewRecorder()
			recorder.Header().Set("Content-Type", "text/plain")
			recorder.WriteHeader(200)
			_, err := recorder.Write([]byte("Hello, World!"))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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
			if writeErr != nil {
				t.Fatalf("unexpected error: %v", writeErr)
			}

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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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

			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if len(readerCalls) <= 1 {
				t.Errorf("expected more than 1 read call, got %d", len(readerCalls))
			}
			if len(closeCalls) != 1 {
				t.Errorf("expected 1 close call, got %d", len(closeCalls))
			}
			if readerCalls[len(readerCalls)-1].ResultErr != errAtEOF {
				t.Errorf("expected error %v, got %v", errAtEOF, readerCalls[len(readerCalls)-1].ResultErr)
			}
			if closeCalls[0].ResultErr != cErr {
				t.Errorf("expected error %v, got %v", cErr, closeCalls[0].ResultErr)
			}

			bodyContent, bodyErr := io.ReadAll(resp.Body)
			if bodyErr != errAtEOF {
				t.Errorf("expected error %v, got %v", errAtEOF, bodyErr)
			}
			if string(bodyContent) != "Hello, World!" {
				t.Errorf("expected %q, got %q", "Hello, World!", string(bodyContent))
			}
			if err := resp.Body.Close(); err != cErr {
				t.Errorf("expected error %v, got %v", cErr, err)
			}
		})
	})
}
