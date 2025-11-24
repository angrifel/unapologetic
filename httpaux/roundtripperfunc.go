package httpaux

import "net/http"

var _ http.RoundTripper = RoundTripperFunc(nil)

// RoundTripperFunc is an adapter type to allow the use of a function as an http.RoundTripper.
// It implements the RoundTrip method by invoking the function itself, passing the request as an argument.
type RoundTripperFunc func(tereq *http.Request) (*http.Response, error)

// RoundTrip executes the HTTP request using the function underlying the RoundTripperFunc and returns the response or error.
func (f RoundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
