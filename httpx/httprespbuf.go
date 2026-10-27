package httpx

import (
	"net/http"

	"github.com/angrifel/unapologetic/iox"
)

// BufferResponseBody returns a copy of resp whose body is wrapped with [iox.ReadSeekCloser].
// The original body is read on demand and its content is buffered in memory as it is read,
// so bytes already read can be read again by seeking. The returned body implements [io.Seeker]
// and can be type-asserted to it. Seek works within the bytes read so far; see [iox.ReadSeekCloser].
// Closing the returned body closes the original body.
// Errors returned by the original body's Read and Close methods are returned unchanged by the
// returned body's Read and Close methods.
func BufferResponseBody(resp *http.Response) *http.Response {
	newBody := iox.ReadSeekCloser(resp.Body)
	result := CloneHTTPResponseWithBody(resp, newBody)

	return result
}
