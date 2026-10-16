package httpaux

import (
	"net/http"

	"github.com/angrifel/unapologetic/ioaux"
)

// BufferResponseBody returns a copy of resp whose body is wrapped with [ioaux.ReadSeekCloser].
// The original body is read on demand and its content is buffered in memory as it is read,
// so bytes already read can be read again by seeking. The returned body implements [io.Seeker]
// and can be type-asserted to it. Closing the returned body closes the original body.
// Errors returned by the original body's Read and Close methods are returned unchanged by the
// returned body's Read and Close methods.
func BufferResponseBody(resp *http.Response) *http.Response {
	newBody := ioaux.ReadSeekCloser(resp.Body)
	result := CloneHTTPResponseWithBody(resp, newBody)

	return result
}
