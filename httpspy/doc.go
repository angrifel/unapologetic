// Package httpspy provides utilities for capturing and inspecting HTTP request and response messages.
//
// This package offers configurable capture functions that can snapshot HTTP messages
// with support for selective header censoring and body collection. It is designed for
// debugging, logging, and testing scenarios where you need to capture HTTP traffic
// without permanently modifying the original message headers.
//
// # Key Features
//
//   - Capture incoming HTTP requests (http.Request) using RequestCaptureWithHeaderCensorshipFunction
//   - Capture outgoing HTTP requests with RequestOutCaptureWithHeaderCensorshipFunction
//   - Capture HTTP responses (http.Response) using ResponseCaptureWithHeaderCensorshipFunction
//   - Selectively censor sensitive headers (e.g., Authorization, Set-Cookie)
//   - Optionally capture message bodies
//   - Non-destructive: headers are temporarily modified during capture then restored
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
package httpspy
