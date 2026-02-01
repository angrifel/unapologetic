package httpaux

import (
	"net/http"
	"time"
)

const defaultTracerTimeout = 30 * time.Second

type (
	roundTripperOptions struct {
		innerRoundTripper          http.RoundTripper
		tracerTimeout              time.Duration
		responseBufferingEnabled   bool
		captureRequest             bool
		requestAttemptCollectBody  bool
		requestCensorHeaders       []string
		requestHeaderCensorText    string
		captureResponse            bool
		responseAttemptCollectBody bool
		responseCensorHeaders      []string
		responseHeaderCensorText   string
		recordTelemetry            TelemetryRecorderFunc
	}

	// RoundTripperOption is an interface for applying custom configurations to.
	RoundTripperOption interface {
		apply(opts *roundTripperOptions)
	}
	funcRoundTripperOption func(*roundTripperOptions)
)

func (f funcRoundTripperOption) apply(opts *roundTripperOptions) { f(opts) }

func WithInnerRoundTripper(inner http.RoundTripper) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.innerRoundTripper = inner })
}

func WithTracerTimeout(timeout time.Duration) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.tracerTimeout = timeout })
}

func WithTelemetryRecorder(recorder TelemetryRecorderFunc) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.recordTelemetry = recorder })
}

func WithResponseBufferingEnabled(enabled bool) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.responseBufferingEnabled = enabled })
}

func WithCaptureRequest(enabled bool) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.captureRequest = enabled })
}

func WithCaptureResponse(enabled bool) RoundTripperOption {
	return funcRoundTripperOption(func(opts *roundTripperOptions) { opts.captureResponse = enabled })
}

func defaultRoundTripperOptions() roundTripperOptions {
	return roundTripperOptions{ //nolint:exhaustruct //it is ok to not use exhaustive struct initialization
		innerRoundTripper: http.DefaultTransport,
		tracerTimeout:     defaultTracerTimeout,
		recordTelemetry:   defaultTelemetryRecorderFunc,
	}
}

func sanitizeRoundTripperOptions(opts *roundTripperOptions) {
	if opts.innerRoundTripper == nil {
		opts.innerRoundTripper = http.DefaultTransport
	}

	if opts.tracerTimeout == 0 {
		opts.tracerTimeout = defaultTracerTimeout
	}

	if opts.recordTelemetry == nil {
		opts.recordTelemetry = defaultTelemetryRecorderFunc
	}
}
