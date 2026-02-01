package httpspy

import (
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"slices"
)

// RequestCaptureWithHeaderCensorshipFunction creates a function that captures incoming HTTP requests.
// If headersToCensor is specified, the returned function will temporarily replace
// those header values with censorText during capture.
func RequestCaptureWithHeaderCensorshipFunction(attemptCollectBody bool, headersToCensor []string, censorText string) func(*http.Request) ([]byte, error) {
	if len(headersToCensor) == 0 {
		return func(request *http.Request) ([]byte, error) {
			return httputil.DumpRequest(request, attemptCollectBody)
		}
	}

	return func(request *http.Request) ([]byte, error) {
		restoreHeadersFn := censorHeaders(request.Header, headersToCensor, censorText)
		defer restoreHeadersFn()

		return httputil.DumpRequest(request, attemptCollectBody)
	}
}

// RequestOutCaptureWithHeaderCensorshipFunction creates a function that captures outgoing HTTP requests.
// If headersToCensor is specified, the returned function will temporarily replace
// those header values with censorText during capture.
func RequestOutCaptureWithHeaderCensorshipFunction(attemptCollectBody bool, headersToCensor []string, censorText string) func(*http.Request) ([]byte, error) {
	if len(headersToCensor) == 0 {
		return func(request *http.Request) ([]byte, error) {
			return httputil.DumpRequestOut(request, attemptCollectBody)
		}
	}

	return func(request *http.Request) ([]byte, error) {
		restoreHeadersFn := censorHeaders(request.Header, headersToCensor, censorText)
		defer restoreHeadersFn()

		return httputil.DumpRequestOut(request, attemptCollectBody)
	}
}

// ResponseCaptureWithHeaderCensorshipFunction creates a function that captures HTTP responses.
// If headersToCensor is specified, the returned function will temporarily replace
// those header values with censorText during capture.
func ResponseCaptureWithHeaderCensorshipFunction(attemptCollectBody bool, headersToCensor []string, censorText string) func(*http.Response) ([]byte, error) {
	if len(headersToCensor) == 0 {
		return func(response *http.Response) ([]byte, error) {
			return httputil.DumpResponse(response, attemptCollectBody)
		}
	}

	return func(response *http.Response) ([]byte, error) {
		restoreHeadersFn := censorHeaders(response.Header, headersToCensor, censorText)
		defer restoreHeadersFn()

		return httputil.DumpResponse(response, attemptCollectBody)
	}
}

func censorHeaders(header http.Header, censorHeader []string, censorText string) (restoreHeaders func()) {
	var preservedHeaders = make(http.Header)

	restoreHeaders = func() {
		for h, hv := range preservedHeaders {
			header[h] = hv
		}
	}

	for _, h := range censorHeader {
		ch := textproto.CanonicalMIMEHeaderKey(h)
		if hv, exists := header[ch]; exists {
			preservedHeaders[ch] = hv
			header[ch] = slices.Repeat([]string{censorText}, len(hv))
		}
	}

	return restoreHeaders
}
