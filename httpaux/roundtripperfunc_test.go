package httpaux

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
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
			if expectedRequest != req {
				t.Errorf("expected request %v, got %v", expectedRequest, req)
			}
			return expectedResponse, nil
		})

		// Act
		actualResponse, err := rt.RoundTrip(expectedRequest)

		// Assert
		if !called {
			t.Error("expected underlying function to be called")
		}
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if expectedResponse != actualResponse {
			t.Errorf("expected response %v, got %v", expectedResponse, actualResponse)
		}
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
		if err != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if actualResponse != nil {
			t.Errorf("expected nil response, got %v", actualResponse)
		}
	})
}
