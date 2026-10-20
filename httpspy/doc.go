// Package httpspy provides utilities for capturing and inspecting HTTP request and response messages,
// and for tracing HTTP round trips.
//
// This package offers configurable capture functions that can snapshot HTTP messages
// with support for selective header censoring and body collection. It is designed for
// debugging, logging, and testing scenarios where you need to capture HTTP traffic
// without permanently modifying the original message headers.
//
// It also offers RoundTripTracer, which records timing and status metrics for a single
// HTTP round trip.
//
// # Key Features
//
//   - Capture incoming HTTP requests (http.Request) using RequestCaptureWithHeaderCensorshipFunction
//   - Capture outgoing HTTP requests with RequestOutCaptureWithHeaderCensorshipFunction
//   - Capture HTTP responses (http.Response) using ResponseCaptureWithHeaderCensorshipFunction
//   - Selectively censor sensitive headers (e.g., Authorization, Set-Cookie)
//   - Optionally capture message bodies
//   - Non-destructive: headers are temporarily modified during capture then restored
//   - Trace HTTP round trips with connection timing metrics using RoundTripTracer
//
// # Basic Usage
//
// Generate capture functions with the desired options:
//
//	// For incoming requests (e.g., in HTTP handlers)
//	captureRequest := httpspy.RequestCaptureWithHeaderCensorshipFunction(
//	    true,                                        // attemptCollectBody
//	    []string{"Authorization", "Cookie"},         // headersToCensor
//	    "*** REDACTED ***",                          // censorText
//	)
//	snapshot, err := captureRequest(req)
//
//	// For outgoing requests (e.g., in HTTP clients)
//	captureRequestOut := httpspy.RequestOutCaptureWithHeaderCensorshipFunction(
//	    true,                                        // attemptCollectBody
//	    []string{"Authorization", "Cookie"},         // headersToCensor
//	    "*** REDACTED ***",                          // censorText
//	)
//	snapshot, err := captureRequestOut(req)
//
//	// For responses
//	captureResponse := httpspy.ResponseCaptureWithHeaderCensorshipFunction(
//	    true,                                        // attemptCollectBody
//	    []string{"Set-Cookie"},                      // headersToCensor
//	    "*** REDACTED ***",                          // censorText
//	)
//	snapshot, err := captureResponse(resp)
//
// # Header Censoring
//
// When headersToCensor is specified, the capture functions will temporarily replace
// the values of specified headers with censorText during the snapshot operation,
// then restore the original values. This ensures sensitive data in the headers is not
// included in captured snapshots while keeping the actual HTTP message intact.
//
// Header names are case-insensitive (following HTTP standards).
//
// # Round-Trip Tracing
//
// RoundTripTracer records connection timing (DNS lookup, dial, TLS handshake), send and
// receive timing, the status code, response body bytes read, and any errors for a single
// HTTP round trip. Instrument the request with TraceRequest and the response with
// TraceResponse, then read the collected metrics with Status:
//
//	tracer := httpspy.NewRoundTripTracer()
//	req = tracer.TraceRequest(req)
//	resp, err := http.DefaultTransport.RoundTrip(req)
//	if err != nil {
//	    tracer.TraceRoundTripError(err)
//	} else {
//	    resp = tracer.TraceResponse(resp)
//	    // read and close resp.Body
//	}
//	status := tracer.Status()
//
// Timing is measured from the GetConn event, which net/http's Transport fires when it begins
// acquiring a connection. Round trips without it, such as a response served by
// http.NewFileTransport or an error before a connection was requested, report zero for Start
// and every timing field.
//
// NewRoundTripTracerWithNotification additionally sends a status snapshot to a channel
// whenever an event in the given TraceEvent mask occurs. Events can be combined with
// bitwise OR (e.g., DNSStart|DNSDone). Notifications are sent synchronously and block until
// received, so the channel must be read concurrently or be buffered.
//
// A RoundTripTracer traces a single round trip: one call to the RoundTrip method of an
// http.RoundTripper. It is not a tracer for http.Client, whose Do method can make several round
// trips for one request (to follow redirects or to retry) that would all be recorded into the
// same tracer. Call RoundTrip directly, and use a new RoundTripTracer for each round trip.
package httpspy
