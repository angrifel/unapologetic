package httpspy

import (
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"slices"
)

// MessageCaptureConfiguration configures how HTTP messages are captured and stored.
type MessageCaptureConfiguration struct {
	// Enabled determines whether message capturing is active.
	Enabled bool
	// AttemptCollectBody determines whether to capture the message body.
	AttemptCollectBody bool
	// CensorHeaders is a list of header names whose values should be censored.
	CensorHeaders []string
	// HeaderCensorText is the replacement text for censored header values.
	HeaderCensorText string
}

// NewRequestCaptureFunction creates a function that captures incoming HTTP requests.
// If capturing is disabled, returns a function that returns empty snapshots.
// If CensorHeaders is specified, the returned function will temporarily replace
// those header values with HeaderCensorText during capture.
func NewRequestCaptureFunction(conf MessageCaptureConfiguration) func(*http.Request) ([]byte, error) {
	if !conf.Enabled {
		return nilRequestCaptureFunction
	}

	if len(conf.CensorHeaders) == 0 {
		return func(request *http.Request) ([]byte, error) {
			return httputil.DumpRequest(request, conf.AttemptCollectBody)
		}
	}

	return func(request *http.Request) ([]byte, error) {
		restoreHeadersFn := censorHeaders(request.Header, conf.CensorHeaders, conf.HeaderCensorText)
		defer restoreHeadersFn()

		return httputil.DumpRequest(request, conf.AttemptCollectBody)
	}
}

// NewRequestOutCaptureFunction creates a function that captures outgoing HTTP requests.
// If capturing is disabled, returns a function that returns empty snapshots.
// If CensorHeaders is specified, the returned function will temporarily replace
// those header values with HeaderCensorText during capture.
func NewRequestOutCaptureFunction(conf MessageCaptureConfiguration) func(*http.Request) ([]byte, error) {
	if !conf.Enabled {
		return nilRequestCaptureFunction
	}

	if len(conf.CensorHeaders) == 0 {
		return func(request *http.Request) ([]byte, error) {
			return httputil.DumpRequestOut(request, conf.AttemptCollectBody)
		}
	}

	return func(request *http.Request) ([]byte, error) {
		restoreHeadersFn := censorHeaders(request.Header, conf.CensorHeaders, conf.HeaderCensorText)
		defer restoreHeadersFn()

		return httputil.DumpRequestOut(request, conf.AttemptCollectBody)
	}
}

// NewResponseCaptureFunction creates a function that captures HTTP responses.
// If capturing is disabled, returns a function that returns empty snapshots.
// If CensorHeaders is specified, the returned function will temporarily replace
// those header values with HeaderCensorText during capture.
func NewResponseCaptureFunction(conf MessageCaptureConfiguration) func(*http.Response) ([]byte, error) {
	if !conf.Enabled {
		return nilResponseCaptureFunction
	}

	if len(conf.CensorHeaders) == 0 {
		return func(response *http.Response) ([]byte, error) {
			return httputil.DumpResponse(response, conf.AttemptCollectBody)
		}
	}

	return func(response *http.Response) ([]byte, error) {
		restoreHeadersFn := censorHeaders(response.Header, conf.CensorHeaders, conf.HeaderCensorText)
		defer restoreHeadersFn()

		return httputil.DumpResponse(response, conf.AttemptCollectBody)
	}
}

func nilRequestCaptureFunction(_ *http.Request) ([]byte, error) { return nil, nil }

func nilResponseCaptureFunction(_ *http.Response) ([]byte, error) { return nil, nil }

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
