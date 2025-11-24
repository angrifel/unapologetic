package httpaux

import (
	"io"
	"net/http"
	"slices"
)

// CloneHTTPResponseWithBody creates a deep copy of an http.Response object, replacing its body with the provided io.ReadCloser.
func CloneHTTPResponseWithBody(response *http.Response, body io.ReadCloser) *http.Response {
	result := *response
	result.TransferEncoding = slices.Clone(response.TransferEncoding)
	result.Header = response.Header.Clone()
	result.Trailer = response.Trailer.Clone()
	result.Body = body

	return &result
}
