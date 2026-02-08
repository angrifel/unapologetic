package httpaux

import (
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/testaux"
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

	assert.EqualFunc(t, response.Header, clonedResponse.Header, testaux.HeaderEqual)
	assert.EqualFunc(t, response.Trailer, clonedResponse.Trailer, testaux.HeaderEqual)
	assert.EqualFunc(t, response.TransferEncoding, clonedResponse.TransferEncoding, slices.Equal)

	assert.Equal(t, replacementBody, clonedResponse.Body)
}
