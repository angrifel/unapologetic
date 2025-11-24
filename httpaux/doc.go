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
// BufferResponseBody reads the entire response body into memory and replaces it with
// a seekable, reusable reader. This is useful when you need to read the body multiple times:
//
//	resp, _ := http.Get(url)
//	buffered := httpaux.BufferResponseBody(resp)
//	// Can now read buffered.Body multiple times
//
// The original body is properly closed, and any read or close errors are preserved
// and returned by the buffered body.
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
