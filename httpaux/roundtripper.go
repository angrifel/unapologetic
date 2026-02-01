package httpaux

import (
	"context"
	"net/http"
	"time"

	"github.com/angrifel/unapologetic/httpspy"
)

type RoundTripTracerStatus = httpspy.RoundTripTracerStatus

type CapturedMessages struct {
	RequestBytes     []byte
	RequestBytesErr  error
	ResponseBytes    []byte
	ResponseBytesErr error
}

// TelemetryRecorderFunc is a function type for recording telemetry data during HTTP request/response lifecycle.
type TelemetryRecorderFunc func(req *http.Request, status RoundTripTracerStatus, capturesMessages CapturedMessages)

type roundTripper struct {
	roundTripperOptions
	requestCapture  func(req *http.Request) ([]byte, error)
	responseCapture func(resp *http.Response) ([]byte, error)
	bufferResponse  func(resp *http.Response) *http.Response
}

// NewRoundTripper creates an http.RoundTripper with customizable options for capturing and buffering HTTP messages.
func NewRoundTripper(opts ...RoundTripperOption) http.RoundTripper {
	sanitizedOptions := defaultRoundTripperOptions()
	for _, opt := range opts {
		opt.apply(&sanitizedOptions)
	}

	sanitizeRoundTripperOptions(&sanitizedOptions)

	result := &roundTripper{
		roundTripperOptions: sanitizedOptions,
		requestCapture:      noRequestCapture,
		responseCapture:     noResponseCapture,
		bufferResponse:      noBufferResponseBody,
	}

	if result.responseBufferingEnabled {
		result.bufferResponse = BufferResponseBody
	}

	if result.captureRequest {
		result.requestCapture = httpspy.RequestOutCaptureWithHeaderCensorshipFunction(sanitizedOptions.requestAttemptCollectBody, sanitizedOptions.requestCensorHeaders, sanitizedOptions.requestHeaderCensorText)
	}

	if result.captureResponse {
		result.responseCapture = httpspy.ResponseCaptureWithHeaderCensorshipFunction(sanitizedOptions.responseAttemptCollectBody, sanitizedOptions.responseCensorHeaders, sanitizedOptions.responseHeaderCensorText)
	}

	return result
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var captureMessages CapturedMessages

	tracer, onExit := rt.newTracer(rt.tracerTimeout)
	nonCancelableCtx := context.WithoutCancel(req.Context())

	defer onExit(req.WithContext(nonCancelableCtx), captureMessages)

	captureMessages.RequestBytes, captureMessages.RequestBytesErr = rt.requestCapture(req)
	req = tracer.TraceRequest(req)

	resp, respErr := rt.innerRoundTripper.RoundTrip(req)
	if respErr != nil {
		tracer.TraceRoundTripError(respErr)

		return nil, respErr
	}

	resp = tracer.TraceResponse(resp)
	resp = rt.bufferResponse(resp)
	captureMessages.ResponseBytes, captureMessages.ResponseBytesErr = rt.responseCapture(resp)

	return resp, nil
}

func (rt *roundTripper) newTracer(timeout time.Duration) (*httpspy.RoundTripTracer, func(req *http.Request, capturesMessages CapturedMessages)) {
	notificationsChan := make(chan RoundTripTracerStatus, 1)
	telemetryTimeout := make(chan struct{})
	tracer := httpspy.NewRoundTripTracerWithNotification(httpspy.BodyClose|httpspy.RoundTripError, notificationsChan)
	onExit := func(req *http.Request, capturesMessages CapturedMessages) {
		go func() {
			defer close(telemetryTimeout)

			time.Sleep(timeout)
		}()

		go func() {
			defer close(notificationsChan)

			select {
			case status := <-notificationsChan:
				rt.recordTelemetry(req, status, capturesMessages)
			case <-telemetryTimeout:
				rt.recordTelemetry(req, tracer.Status(), capturesMessages)
			}
		}()
	}

	return tracer, onExit
}

func noBufferResponseBody(resp *http.Response) *http.Response { return resp }
func noRequestCapture(_ *http.Request) ([]byte, error)        { return nil, nil }
func noResponseCapture(_ *http.Response) ([]byte, error)      { return nil, nil }

func defaultTelemetryRecorderFunc(_ *http.Request, _ RoundTripTracerStatus, _ CapturedMessages) {}
