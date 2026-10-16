// Package httpaux provides utilities for working with HTTP responses and round trippers.
//
// This package offers helper functions for common HTTP client operations:
//   - Cloning http.Response objects with custom bodies
//   - Buffering response bodies into memory for multiple reads
//   - Creating http.RoundTripper implementations from functions
//
// # Response Cloning
//
// CloneHTTPResponseWithBody creates a deep copy of an http.Response while replacing
// the body with a new io.ReadCloser:
//
//	newBody := io.NopCloser(bytes.NewReader(data))
//	cloned := httpaux.CloneHTTPResponseWithBody(originalResp, newBody)
//
// # Response Buffering
//
// BufferResponseBody replaces the response body with a seekable reader that buffers content
// in memory as it is read. This is useful when you need to read the body multiple times:
//
//	resp, _ := http.Get(url)
//	buffered := httpaux.BufferResponseBody(resp)
//	data, _ := io.ReadAll(buffered.Body)
//	_, _ = buffered.Body.(io.Seeker).Seek(0, io.SeekStart)
//	// buffered.Body can now be read again
//
// The original body is read on demand, not up front, and is closed when the buffered body
// is closed. Any read or close errors from the original body are returned by the buffered body.
//
// # RoundTripper Adapter
//
// RoundTripperFunc allows using a function as an http.RoundTripper, similar to
// how http.HandlerFunc works for handlers:
//
//	rt := httpaux.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
//	    // Custom request handling
//	    return &http.Response{StatusCode: 200}, nil
//	})
//	client := &http.Client{Transport: rt}
package httpaux
