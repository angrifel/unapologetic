// Package httpspy provides utilities for capturing and inspecting HTTP request and response messages.
//
// This package offers configurable capture functions that can snapshot HTTP messages
// with support for selective header censoring and body collection. It is designed for
// debugging, logging, and testing scenarios where you need to capture HTTP traffic
// without permanently modifying the original message headers.
//
// # Key Features
//
//   - Capture incoming HTTP requests (http.Request) using NewRequestCaptureFunction
//   - Capture outgoing HTTP requests with NewRequestOutCaptureFunction
//   - Capture HTTP responses (http.Response) using NewResponseCaptureFunction
//   - Selectively censor sensitive headers (e.g., Authorization, Set-Cookie)
//   - Optionally capture message bodies
//   - Non-destructive: headers are temporarily modified during capture then restored
//
// # Basic Usage
//
// Create a capture configuration and generate capture functions:
//
//	config := httpspy.MessageCaptureConfiguration{
//	    Enabled:            true,
//	    AttemptCollectBody: true,
//	    CensorHeaders:      []string{"Authorization", "Cookie"},
//	    HeaderCensorText:   "*** REDACTED ***",
//	}
//
//	// For incoming requests (e.g., in HTTP handlers)
//	captureRequest := httpspy.NewRequestCaptureFunction(config)
//	snapshot, err := captureRequest(req)
//
//	// For outgoing requests (e.g., in HTTP clients)
//	captureRequestOut := httpspy.NewRequestOutCaptureFunction(config)
//	snapshot, err := captureRequestOut(req)
//
//	// For responses
//	captureResponse := httpspy.NewResponseCaptureFunction(config)
//	snapshot, err := captureResponse(resp)
//
// # Header Censoring
//
// When CensorHeaders is specified, the capture functions will temporarily replace
// the values of specified headers with HeaderCensorText during the snapshot operation,
// then restore the original values. This ensures sensitive data is not included in
// captured snapshots while keeping the actual HTTP message intact.
//
// Header names are case-insensitive (following HTTP standards).
//
// # Disabling Capture
//
// When Enabled is false, the capture functions return nil snapshots without any
// overhead, making it easy to toggle capture on/off at runtime.
package httpspy
