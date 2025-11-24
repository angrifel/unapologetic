package httpaux

import (
	"net/http"

	"github.com/angrifel/unapologetic/ioaux"
)

// BufferResponseBody drains and closes the response body. The response body is replaced with an in-memory io.ReadCloser
// that returns the content of the response.
// Any errors occurring while reading or closing the original body are
// retained are returned back by the in-memory io.ReadCloser.
func BufferResponseBody(resp *http.Response) *http.Response {
	if resp == nil {
		return nil
	}

	newBody := ioaux.ReadSeekCloser(resp.Body)
	result := CloneHTTPResponseWithBody(resp, newBody)

	return result
}
