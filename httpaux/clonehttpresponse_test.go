package httpaux

import (
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestCloneHTTPResponseWithSameBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	http.Error(recorder, "Internal Server Error", http.StatusInternalServerError)
	response := recorder.Result()
	clonedResponse := CloneHTTPResponseWithBody(response, response.Body)

	assertCloneHTTPResponseWithBody(t, response, response.Body, clonedResponse)
}

func TestCloneHTTPResponseWithDifferentBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	http.Error(recorder, "Internal Server Error", http.StatusInternalServerError)
	response := recorder.Result()
	replacementBody := io.NopCloser(strings.NewReader("internal server error"))

	clonedResponse := CloneHTTPResponseWithBody(response, replacementBody)

	assertCloneHTTPResponseWithBody(t, response, replacementBody, clonedResponse)
}

func assertCloneHTTPResponseWithBody(t *testing.T, response *http.Response, replacementBody io.ReadCloser, clonedResponse *http.Response) {
	if response == clonedResponse {
		t.Error("expected cloned response to be a different instance")
	}
	if replacementBody != response.Body {
		if response.Body == clonedResponse.Body {
			t.Error("expected cloned response body to be a different instance")
		}
	} else {
		if response.Body != clonedResponse.Body {
			t.Error("expected cloned response body to be the same instance")
		}
	}
	if response.Status != clonedResponse.Status {
		t.Errorf("Status: expected %v, got %v", response.Status, clonedResponse.Status)
	}
	if response.StatusCode != clonedResponse.StatusCode {
		t.Errorf("StatusCode: expected %v, got %v", response.StatusCode, clonedResponse.StatusCode)
	}
	if response.Proto != clonedResponse.Proto {
		t.Errorf("Proto: expected %v, got %v", response.Proto, clonedResponse.Proto)
	}
	if response.ProtoMajor != clonedResponse.ProtoMajor {
		t.Errorf("ProtoMajor: expected %v, got %v", response.ProtoMajor, clonedResponse.ProtoMajor)
	}
	if response.ProtoMinor != clonedResponse.ProtoMinor {
		t.Errorf("ProtoMinor: expected %v, got %v", response.ProtoMinor, clonedResponse.ProtoMinor)
	}
	if response.ContentLength != clonedResponse.ContentLength {
		t.Errorf("ContentLength: expected %v, got %v", response.ContentLength, clonedResponse.ContentLength)
	}
	if response.Close != clonedResponse.Close {
		t.Errorf("Close: expected %v, got %v", response.Close, clonedResponse.Close)
	}
	if response.Uncompressed != clonedResponse.Uncompressed {
		t.Errorf("Uncompressed: expected %v, got %v", response.Uncompressed, clonedResponse.Uncompressed)
	}
	if response.Request != clonedResponse.Request {
		t.Errorf("Request: expected %v, got %v", response.Request, clonedResponse.Request)
	}
	if response.TLS != clonedResponse.TLS {
		t.Errorf("TLS: expected %v, got %v", response.TLS, clonedResponse.TLS)
	}
	if !maps.EqualFunc(response.Header, clonedResponse.Header, slices.Equal) {
		t.Errorf("Header: expected %v, got %v", response.Header, clonedResponse.Header)
	}
	if !maps.EqualFunc(response.Trailer, clonedResponse.Trailer, slices.Equal) {
		t.Errorf("Trailer: expected %v, got %v", response.Trailer, clonedResponse.Trailer)
	}
	if !slices.Equal(response.TransferEncoding, clonedResponse.TransferEncoding) {
		t.Errorf("TransferEncoding: expected %v, got %v", response.TransferEncoding, clonedResponse.TransferEncoding)
	}
	if replacementBody != clonedResponse.Body {
		t.Errorf("Body: expected replacement body, got %v", clonedResponse.Body)
	}
}
