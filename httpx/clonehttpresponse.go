package httpx

import (
	"io"
	"net/http"
	"slices"
)

// CloneHTTPResponseWithBody creates a copy of an http.Response object with its Body replaced by the provided io.ReadCloser.
// All fields are cloned except Request and TLS, which are shared with the original response.
func CloneHTTPResponseWithBody(response *http.Response, body io.ReadCloser) *http.Response {
	result := *response
	result.TransferEncoding = slices.Clone(response.TransferEncoding)
	result.Header = response.Header.Clone()
	result.Trailer = response.Trailer.Clone()
	result.Body = body

	return &result
}
