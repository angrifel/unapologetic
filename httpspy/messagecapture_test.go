package httpspy

import (
	"bytes"
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/testaux"
)

//go:embed testmaterial/*
var testMaterial embed.FS

func TestRequestCaptureWithHeaderCensorshipFunction(t *testing.T) {
	type testCase struct {
		testFile           string
		attemptCollectBody bool
		censorHeaders      []string
		headerCensorText   string
		req                *http.Request
	}

	newTestRequest := func(method, target string, body io.Reader, reqUpdate func(r *http.Request)) *http.Request {
		request := httptest.NewRequestWithContext(t.Context(), method, target, body)
		if reqUpdate != nil {
			reqUpdate(request)
		}

		return request
	}

	var testCases = []testCase{
		{
			testFile: "testmaterial/req01",
			req:      newTestRequest("GET", "https://example.com", nil, nil),
		},
		{
			testFile: "testmaterial/req02",
			req: newTestRequest("GET", "https://example.com", nil, func(r *http.Request) {
				r.Header.Set("x-Header-1", "value")
				r.Header.Add("Other-header", "thiswillbecensored")
			}),
			censorHeaders: []string{"other-header"},
		},
		{
			testFile: "testmaterial/req03",
			req: newTestRequest("GET", "https://example.com", nil, func(r *http.Request) {
				r.Header.Set("x-Header-1", "value")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-1")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-2")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-3")
			}),
			censorHeaders:    []string{"other-header-mv"},
			headerCensorText: "*** CENSORED ***",
		},
		{
			testFile: "testmaterial/req04",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
			}),
			attemptCollectBody: true,
		},
		{
			testFile: "testmaterial/req05",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
				r.Header.Set("Authorization", "Bearer secret-token")
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"authorization"},
		},
		{
			testFile: "testmaterial/req06",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
				r.Header.Add("X-Custom-Header", "value1")
				r.Header.Add("X-Custom-Header", "value2")
				r.Header.Add("X-Custom-Header", "value3")
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"x-custom-header"},
			headerCensorText:   "*** CENSORED ***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testFile, func(t *testing.T) {
			// arrange
			tcBodyBytes, tcBodyBytesErr := io.ReadAll(tc.req.Body)
			_ = tcBodyBytesErr
			expectedBytes, expectedBytesErr := testMaterial.ReadFile(tc.testFile)
			if expectedBytesErr != nil {
				t.Fatalf("Failed to read %s: %v", tc.testFile, expectedBytesErr)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// arrange
				headersBeforeCapture := r.Header.Clone()
				captureFunction := RequestCaptureWithHeaderCensorshipFunction(tc.attemptCollectBody, tc.censorHeaders, tc.headerCensorText)

				// act
				requestBytes, requestErr := captureFunction(r)

				// assert
				assert.EqualFunc(t, headersBeforeCapture, r.Header, testaux.HeaderEqual)
				assert.IsNil(t, requestErr)
				assert.Equal(t, string(expectedBytes), string(requestBytes))

				w.WriteHeader(http.StatusOK)
			}))

			t.Cleanup(server.Close)

			host := server.URL[len("http://"):]
			expectedBytes = bytes.ReplaceAll(expectedBytes, []byte("\r\nHost: example.com\r\n"), []byte("\r\nHost: "+host+"\r\n"))
			assert.IsNil(t, testaux.WaitForConnectAccept(host))

			req := httptest.NewRequestWithContext(t.Context(), tc.req.Method, server.URL, bytes.NewReader(tcBodyBytes))
			req.RequestURI = ""
			req.Header = tc.req.Header.Clone()

			// act
			_, respErr := server.Client().Do(req)

			// assert
			assert.IsNil(t, respErr)
		})
	}
}

func TestRequestOutCaptureWithHeaderCensorshipFunction(t *testing.T) {
	type testCase struct {
		testFile           string
		attemptCollectBody bool
		censorHeaders      []string
		headerCensorText   string
		req                *http.Request
	}

	newTestRequest := func(method, target string, body io.Reader, reqUpdate func(r *http.Request)) *http.Request {
		request := httptest.NewRequestWithContext(t.Context(), method, target, body)
		if reqUpdate != nil {
			reqUpdate(request)
		}

		return request
	}

	var testCases = []testCase{
		{
			testFile: "testmaterial/reqout01",
			req:      newTestRequest("GET", "https://example.com", nil, nil),
		},
		{
			testFile: "testmaterial/reqout02",
			req: newTestRequest("GET", "https://example.com", nil, func(r *http.Request) {
				r.Header.Set("x-Header-1", "value")
				r.Header.Add("Other-header", "thiswillbecensored")
			}),
			censorHeaders: []string{"other-header"},
		},
		{
			testFile: "testmaterial/reqout03",
			req: newTestRequest("GET", "https://example.com", nil, func(r *http.Request) {
				r.Header.Set("x-Header-1", "value")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-1")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-2")
				r.Header.Add("Other-header-Mv", "thiswillbecensored-3")
			}),
			censorHeaders:    []string{"other-header-mv"},
			headerCensorText: "*** CENSORED ***",
		},
		{
			testFile: "testmaterial/reqout04",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
			}),
			attemptCollectBody: true,
		},
		{
			testFile: "testmaterial/reqout05",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
				r.Header.Set("Authorization", "Bearer secret-token")
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"authorization"},
		},
		{
			testFile: "testmaterial/reqout06",
			req: newTestRequest("POST", "https://example.com", strings.NewReader("test request body"), func(r *http.Request) {
				r.Header.Set("Content-Type", "text/plain")
				r.Header.Add("X-Custom-Header", "value1")
				r.Header.Add("X-Custom-Header", "value2")
				r.Header.Add("X-Custom-Header", "value3")
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"x-custom-header"},
			headerCensorText:   "*** CENSORED ***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testFile, func(t *testing.T) {
			// arrange
			expectedBytes, expectedBytesErr := testMaterial.ReadFile(tc.testFile)
			if expectedBytesErr != nil {
				t.Fatalf("Failed to read %s: %v", tc.testFile, expectedBytesErr)
			}

			headersBeforeCapture := tc.req.Header.Clone()
			captureFunction := RequestOutCaptureWithHeaderCensorshipFunction(tc.attemptCollectBody, tc.censorHeaders, tc.headerCensorText)

			// act
			requestBytes, requestErr := captureFunction(tc.req)

			// assert
			assert.EqualFunc(t, headersBeforeCapture, tc.req.Header, testaux.HeaderEqual)
			assert.IsNil(t, requestErr)
			assert.Equal(t, string(expectedBytes), string(requestBytes))
		})
	}
}

func TestResponseCaptureWithHeaderCensorshipFunction(t *testing.T) {
	type testCase struct {
		testFile           string
		attemptCollectBody bool
		censorHeaders      []string
		headerCensorText   string
		resp               *http.Response
	}

	newTestResponse := func(body io.Reader, respUpdate func(r *http.Response)) *http.Response {
		resp := &http.Response{
			StatusCode:    200,
			Status:        "200 OK",
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: 0,
			Body:          io.NopCloser(body),
			Request:       httptest.NewRequest("GET", "https://example.com", nil),
		}

		if respUpdate != nil {
			respUpdate(resp)
		}

		return resp
	}

	var testCases = []testCase{
		{
			testFile: "testmaterial/resp01",
			resp: newTestResponse(http.NoBody, func(resp *http.Response) {
				resp.Header = http.Header{"Content-Type": []string{"text/plain"}}
			}),
		},
		{
			testFile: "testmaterial/resp02",
			resp: newTestResponse(http.NoBody, func(resp *http.Response) {
				resp.Header = http.Header{
					"Content-Type": []string{"text/plain"},
					"Set-Cookie":   []string{"session=abc123; Path=/; HttpOnly"},
				}
			}),
			censorHeaders: []string{"set-cookie"},
		},
		{
			testFile: "testmaterial/resp03",
			resp: newTestResponse(http.NoBody, func(resp *http.Response) {
				resp.Header = http.Header{
					"Content-Type":      []string{"text/plain"},
					"X-Custom-Response": []string{"value1", "value2"},
				}
			}),
			censorHeaders:    []string{"x-custom-response"},
			headerCensorText: "*** CENSORED ***",
		},
		{
			testFile: "testmaterial/resp04",
			resp: newTestResponse(strings.NewReader("test response body"), func(resp *http.Response) {
				resp.ContentLength = 18
				resp.Header = http.Header{"Content-Type": []string{"text/plain"}}
			}),
			attemptCollectBody: true,
		},
		{
			testFile: "testmaterial/resp05",
			resp: newTestResponse(strings.NewReader("test response body"), func(resp *http.Response) {
				resp.Header = http.Header{
					"Content-Type": []string{"text/plain"},
					"Set-Cookie":   []string{"session=abc123; Path=/; HttpOnly"},
				}
				resp.ContentLength = 18
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"set-cookie"},
		},
		{
			testFile: "testmaterial/resp06",
			resp: newTestResponse(strings.NewReader("test response body"), func(resp *http.Response) {
				resp.Header = http.Header{
					"Content-Type":      []string{"text/plain"},
					"X-Custom-Response": []string{"value1", "value2"},
				}
				resp.ContentLength = 18
			}),
			attemptCollectBody: true,
			censorHeaders:      []string{"x-custom-response"},
			headerCensorText:   "*** CENSORED ***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testFile, func(t *testing.T) {
			// arrange
			expectedBytes, expectedBytesErr := testMaterial.ReadFile(tc.testFile)
			if expectedBytesErr != nil {
				t.Fatalf("Failed to read %s: %v", tc.testFile, expectedBytesErr)
			}

			headersBeforeCapture := tc.resp.Header.Clone()
			captureFunction := ResponseCaptureWithHeaderCensorshipFunction(tc.attemptCollectBody, tc.censorHeaders, tc.headerCensorText)

			// act
			responseBytes, requestErr := captureFunction(tc.resp)

			// assert
			assert.EqualFunc(t, headersBeforeCapture, tc.resp.Header, testaux.HeaderEqual)
			assert.IsNil(t, requestErr)
			assert.Equal(t, string(expectedBytes), string(responseBytes))
		})
	}
}
