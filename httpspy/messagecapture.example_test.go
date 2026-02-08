package httpspy_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/angrifel/unapologetic/httpspy"
)

func ExampleRequestCaptureWithHeaderCensorshipFunction() {
	capture := httpspy.RequestCaptureWithHeaderCensorshipFunction(false, []string{"authorization"}, "[REDACTED]")

	req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
	req.Header.Set("Authorization", "Bearer secret-token")

	data, _ := capture(req)
	fmt.Println(strings.Contains(string(data), "[REDACTED]"))
	fmt.Println(!strings.Contains(string(data), "secret-token"))
	// Output:
	// true
	// true
}

func ExampleRequestOutCaptureWithHeaderCensorshipFunction() {
	capture := httpspy.RequestOutCaptureWithHeaderCensorshipFunction(true, nil, "")

	req := httptest.NewRequest(http.MethodPost, "https://example.com", strings.NewReader("request body"))
	req.Header.Set("Content-Type", "text/plain")

	data, _ := capture(req)
	fmt.Println(len(data) > 0)
	// Output:
	// true
}

func ExampleResponseCaptureWithHeaderCensorshipFunction() {
	capture := httpspy.ResponseCaptureWithHeaderCensorshipFunction(true, []string{"set-cookie"}, "[REDACTED]")

	resp := &http.Response{
		StatusCode:    200,
		Status:        "200 OK",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 13,
		Header: http.Header{
			"Content-Type": []string{"text/plain"},
			"Set-Cookie":   []string{"session=abc123"},
		},
		Body:    io.NopCloser(strings.NewReader("response body")),
		Request: httptest.NewRequest("GET", "https://example.com", nil),
	}

	data, _ := capture(resp)
	fmt.Println(strings.Contains(string(data), "[REDACTED]"))
	fmt.Println(!strings.Contains(string(data), "session=abc123"))
	fmt.Println(strings.Contains(string(data), "response body"))
	// Output:
	// true
	// true
	// true
}
