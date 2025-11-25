package httpaux

import (
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
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
	assert.NotEqual(t, response, clonedResponse)
	if replacementBody != response.Body {
		assert.NotEqual(t, response.Body, clonedResponse.Body)
	} else {
		assert.Equal(t, response.Body, clonedResponse.Body)
	}

	assert.Equal(t, response.Status, clonedResponse.Status)
	assert.Equal(t, response.StatusCode, clonedResponse.StatusCode)
	assert.Equal(t, response.Proto, clonedResponse.Proto)
	assert.Equal(t, response.ProtoMajor, clonedResponse.ProtoMajor)
	assert.Equal(t, response.ProtoMinor, clonedResponse.ProtoMinor)
	assert.Equal(t, response.ContentLength, clonedResponse.ContentLength)
	assert.Equal(t, response.Close, clonedResponse.Close)
	assert.Equal(t, response.Uncompressed, clonedResponse.Uncompressed)
	assert.Equal(t, response.Request, clonedResponse.Request)
	assert.Equal(t, response.TLS, clonedResponse.TLS)

	if !maps.EqualFunc(response.Header, clonedResponse.Header, slices.Equal) {
		t.Errorf("Header: expected %v, got %v", response.Header, clonedResponse.Header)
	}
	if !maps.EqualFunc(response.Trailer, clonedResponse.Trailer, slices.Equal) {
		t.Errorf("Trailer: expected %v, got %v", response.Trailer, clonedResponse.Trailer)
	}
	if !slices.Equal(response.TransferEncoding, clonedResponse.TransferEncoding) {
		t.Errorf("TransferEncoding: expected %v, got %v", response.TransferEncoding, clonedResponse.TransferEncoding)
	}

	assert.Equal(t, replacementBody, clonedResponse.Body)
}
