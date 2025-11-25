package httpaux

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestRoundTripperFunc(t *testing.T) {
	t.Run("Forwards the request to the underlying function", func(t *testing.T) {
		// Arrange
		expectedRequest, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("response body")),
		}

		called := false
		rt := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			assert.Equal(t, expectedRequest, req)

			return expectedResponse, nil
		})

		// Act
		actualResponse, err := rt.RoundTrip(expectedRequest)

		// Assert
		assert.Equal(t, true, called)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("Returns error from the underlying function", func(t *testing.T) {
		// Arrange
		expectedRequest, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		expectedError := errors.New("network error")

		rt := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, expectedError
		})

		// Act
		actualResponse, err := rt.RoundTrip(expectedRequest)

		// Assert
		assert.Equal(t, expectedError, err)
		assert.Equal(t, nil, actualResponse)
	})
}
