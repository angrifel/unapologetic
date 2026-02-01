package httpspy

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"

	"github.com/angrifel/unapologetic/ioaux"
)

// TraceEvent represents a specific event that occurs during an HTTP round-trip operation.
// Events can be combined using bitwise OR to create event masks for selective notification.
type TraceEvent uint64

const (
	// GetConn is triggered when the HTTP client begins acquiring a connection.
	GetConn TraceEvent = 1 << iota
	// GotConn is triggered when the HTTP client has successfully acquired a connection.
	GotConn
	// DNSStart is triggered when DNS lookup begins for the target host.
	DNSStart
	// DNSDone is triggered when DNS lookup completes for the target host.
	DNSDone
	// ConnectStart is triggered when the TCP connection establishment begins.
	ConnectStart
	// ConnectDone is triggered when the TCP connection establishment completes.
	ConnectDone
	// TLSHandshakeStart is triggered when the TLS handshake begins (HTTPS only).
	TLSHandshakeStart
	// TLSHandshakeDone is triggered when the TLS handshake completes (HTTPS only).
	TLSHandshakeDone
	// WroteRequest is triggered when the HTTP request has been written to the connection.
	WroteRequest
	// WroteRequestError is triggered when an error occurs while writing the request.
	WroteRequestError
	// GotFirstResponseByte is triggered when the first byte of the response is received.
	GotFirstResponseByte
	// ResponseHeaders is triggered when response headers have been received.
	ResponseHeaders
	// RoundTripError is triggered when an error occurs during the round-trip operation.
	RoundTripError
	// BodyRead is triggered each time data is read from the response body.
	BodyRead
	// BodyReadEOF is triggered when EOF is reached while reading the response body.
	BodyReadEOF
	// BodyReadError is triggered when an error occurs while reading the response body.
	BodyReadError
	// BodyClose is triggered when the response body is closed.
	BodyClose
	// BodyCloseError is triggered when an error occurs while closing the response body.
	BodyCloseError
)

// RoundTripTracer tracks timing and status information for HTTP round-trip operations.
// It instruments HTTP requests and responses to capture detailed metrics about connection
// establishment, data transfer, and lifecycle events.
//
// RoundTripTracer should be used only once per round-trip.
type RoundTripTracer struct {
	start                  atomic.Value
	dnsLookupTime          atomic.Int64
	dialTime               atomic.Int64
	tlsHandshakeTime       atomic.Int64
	connectionTime         atomic.Int64
	sendTime               atomic.Int64
	timeToFirstByteTime    atomic.Int64
	receiveTime            atomic.Int64
	closeTime              atomic.Int64
	responseBodyBytesRead  atomic.Int64
	statusCode             atomic.Int64
	sendError              atomic.Value
	roundTripError         atomic.Value
	responseBodyReadError  atomic.Value
	responseBodyCloseError atomic.Value
	statusBits             atomic.Uint64
	notificationMask       TraceEvent
	notificationChan       chan<- RoundTripTracerStatus
}

// NewRoundTripTracer creates a new RoundTripTracer for tracking HTTP round-trip metrics.
// The returned tracer can be used to instrument HTTP requests via TraceRequest and
// responses via TraceResponse.
func NewRoundTripTracer() *RoundTripTracer {
	return &RoundTripTracer{} //nolint:exhaustruct // RoundTripTracer has many fields, but most are optional
}

// NewRoundTripTracerWithNotification creates a new RoundTripTracer that sends notifications
// to the provided channel when specific trace events occur.
//
// The notificationMask parameter specifies which events should trigger notifications using
// bitwise OR of TraceEvent constants (e.g., DNSStart|DNSDone). When a matching event occurs,
// the current RoundTripTracerStatus is sent to notificationChan.
//
// If notificationChan is nil, no notifications are sent but metrics are still collected.
func NewRoundTripTracerWithNotification(notificationMask TraceEvent, notificationChan chan<- RoundTripTracerStatus) *RoundTripTracer {
	return &RoundTripTracer{ //nolint:exhaustruct // RoundTripTracer has many fields, but most are optional
		notificationMask: notificationMask,
		notificationChan: notificationChan,
	}
}

// RoundTripTracerStatus contains the collected metrics and status information from an HTTP
// round-trip operation. All timing fields are measured from the start of the request.
type RoundTripTracerStatus = struct {
	// Start is the time when the HTTP request operation began
	Start time.Time
	// DNSLookupTime represents the duration taken to perform DNS lookup
	DNSLookupTime time.Duration
	// DialTime represents the duration taken to establish a connection to the server
	DialTime time.Duration
	// TLSHandshakeTime represents the duration taken to complete the TLS handshake
	TLSHandshakeTime time.Duration
	// ConnectionTime represents the duration taken to establish the full connection
	ConnectionTime time.Duration
	// SendTime represents the duration taken to send the request to the server
	SendTime time.Duration
	// SendError captures any error encountered while writing the request.
	SendError error
	// FirstByteTime represents the duration until the first byte of the response is received, measured from the beginning of the request.
	FirstByteTime time.Duration
	// ReceiveTime represents the total duration until the full response body is received, measured from the beginning of the request.
	ReceiveTime time.Duration
	// CloseTime represents the duration until the connection is closed, measured from the beginning of the request.
	CloseTime time.Duration
	// ResponseBodyBytesRead represents the number of bytes read from the response body.
	ResponseBodyBytesRead int
	// StatusCode represents the HTTP status code returned by the server during the round-trip operation.
	StatusCode int
	// RoundTripError represents any error encountered during the round trip of the request-response cycle.
	RoundTripError error
	// ResponseBodyReadError represents an error that occurred while reading the response body during the round-trip.
	ResponseBodyReadError error
	// ResponseBodyCloseError represents an error encountered while closing the response body.
	ResponseBodyCloseError error
}

// TraceRequest instruments an HTTP request to capture detailed timing and lifecycle metrics.
// It returns a new request with an httptrace.ClientTrace attached to its context.
//
// The tracer captures metrics for connection establishment (DNS, TCP dial, TLS handshake),
// request transmission, and response timing. The returned request should be used in place
// of the original request when making the HTTP call.
func (rtt *RoundTripTracer) TraceRequest(req *http.Request) *http.Request {
	var (
		dnsStartTime          atomic.Pointer[time.Time]
		dialStartTime         atomic.Pointer[time.Time]
		tlsHandshakeStartTime atomic.Pointer[time.Time]
	)

	ct := &httptrace.ClientTrace{ //nolint:exhaustruct
		GetConn: func(_ string) {
			now := time.Now()
			rtt.start.CompareAndSwap(nil, now)
			rtt.updateStatusBitsAndNotify(GetConn)
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			now := time.Now()
			rtt.connectionTime.Store(int64(now.Sub(loadValue[time.Time](&rtt.start))))
			rtt.updateStatusBitsAndNotify(GotConn)
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			now := time.Now()
			dnsStartTime.Store(&now)
			rtt.updateStatusBitsAndNotify(DNSStart)
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			now := time.Now()
			rtt.dnsLookupTime.Add(int64(now.Sub(*dnsStartTime.Load())))
			rtt.updateStatusBitsAndNotify(DNSDone)
		},
		ConnectStart: func(_, _ string) {
			now := time.Now()
			dialStartTime.Store(&now)
			rtt.updateStatusBitsAndNotify(ConnectStart)
		},
		ConnectDone: func(_, _ string, _ error) {
			now := time.Now()
			rtt.dialTime.Add(int64(now.Sub(*dialStartTime.Load())))
			rtt.updateStatusBitsAndNotify(ConnectDone)
		},
		TLSHandshakeStart: func() {
			now := time.Now()
			tlsHandshakeStartTime.Store(&now)
			rtt.updateStatusBitsAndNotify(TLSHandshakeStart)
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			now := time.Now()
			rtt.tlsHandshakeTime.Add(int64(now.Sub(*tlsHandshakeStartTime.Load())))
			rtt.updateStatusBitsAndNotify(TLSHandshakeDone)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			now := time.Now()
			rtt.sendTime.Store(int64(now.Sub(loadValue[time.Time](&rtt.start))))
			event := WroteRequest
			if info.Err != nil {
				event |= WroteRequestError
				rtt.sendError.Store(&info.Err)
			}

			rtt.updateStatusBitsAndNotify(event)
		},
		GotFirstResponseByte: func() {
			now := time.Now()
			rtt.timeToFirstByteTime.Store(int64(now.Sub(loadValue[time.Time](&rtt.start))))
			rtt.updateStatusBitsAndNotify(GotFirstResponseByte)
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), ct))
}

// TraceResponse instruments an HTTP response to capture response body reading metrics.
// It wraps the response body with instrumentation that tracks bytes read, read errors,
// and close timing.
//
// The returned response should be used in place of the original response. All operations
// on the response body (Read, Close) will be tracked and recorded in the tracer's metrics.
//
// This method should be called immediately after receiving a response from an HTTP call
// that was instrumented with TraceRequest.
func (rtt *RoundTripTracer) TraceResponse(response *http.Response) *http.Response {
	rtt.statusCode.Store(int64(response.StatusCode))
	rtt.updateStatusBitsAndNotify(ResponseHeaders)

	responseBody := response.Body

	response.Body = struct {
		io.Reader
		io.Closer
	}{
		Reader: ioaux.ReaderFunc(func(p []byte) (n int, err error) {
			n, err = responseBody.Read(p)
			now := time.Now()
			rtt.receiveTime.Store(int64(now.Sub(loadValue[time.Time](&rtt.start))))
			rtt.responseBodyBytesRead.Add(int64(n))

			event := BodyRead
			if err != nil {
				if err == io.EOF {
					event |= BodyReadEOF
				} else {
					event |= BodyReadError
					rtt.responseBodyReadError.Store(&err)
				}
			}

			rtt.updateStatusBitsAndNotify(event)

			return
		}),
		Closer: ioaux.CloserFunc(func() (err error) {
			err = responseBody.Close()
			now := time.Now()
			rtt.closeTime.Store(int64(now.Sub(loadValue[time.Time](&rtt.start))))

			event := BodyClose
			if err != nil {
				event |= BodyCloseError
				rtt.responseBodyCloseError.Store(&err)
			}

			rtt.updateStatusBitsAndNotify(event)

			return
		}),
	}

	return response
}

// TraceRoundTripError records an error that occurred during the HTTP round-trip operation.
// This method should be called when an error prevents the normal completion of the request,
// such as connection failures, timeouts, or other transport-level errors.
//
// The error is stored and timing fields (ReceiveTime, CloseTime) are set to mark when
// the error occurred. The RoundTripError status bit is set and notifications are sent
// if configured.
func (rtt *RoundTripTracer) TraceRoundTripError(err error) {
	now := time.Now()
	start := loadValue[time.Time](&rtt.start)

	rtt.roundTripError.Store(&err)
	rtt.receiveTime.Store(int64(now.Sub(start)))
	rtt.closeTime.Store(int64(now.Sub(start)))
	rtt.updateStatusBitsAndNotify(RoundTripError)
}

// Status returns a snapshot of the current tracer metrics and status.
// This method can be called at any time during or after the HTTP round-trip to retrieve
// the accumulated timing data, status code, errors, and other metrics.
//
// The returned RoundTripTracerStatus contains all collected metrics up to the point of
// the call. For ongoing operations, some fields may be zero until the corresponding
// events occur.
func (rtt *RoundTripTracer) Status() RoundTripTracerStatus {
	return RoundTripTracerStatus{
		Start:                  loadValue[time.Time](&rtt.start),
		DNSLookupTime:          time.Duration(rtt.dnsLookupTime.Load()),
		DialTime:               time.Duration(rtt.dialTime.Load()),
		TLSHandshakeTime:       time.Duration(rtt.tlsHandshakeTime.Load()),
		ConnectionTime:         time.Duration(rtt.connectionTime.Load()),
		SendTime:               time.Duration(rtt.sendTime.Load()),
		SendError:              loadValue[error](&rtt.sendError),
		FirstByteTime:          time.Duration(rtt.timeToFirstByteTime.Load()),
		ResponseBodyBytesRead:  int(rtt.responseBodyBytesRead.Load()),
		ReceiveTime:            time.Duration(rtt.receiveTime.Load()),
		CloseTime:              time.Duration(rtt.closeTime.Load()),
		StatusCode:             int(rtt.statusCode.Load()),
		RoundTripError:         loadValue[error](&rtt.roundTripError),
		ResponseBodyReadError:  loadValue[error](&rtt.responseBodyReadError),
		ResponseBodyCloseError: loadValue[error](&rtt.responseBodyCloseError),
	}
}

func loadValue[T any](v *atomic.Value) T {
	if value := v.Load(); value != nil {
		return value.(T) //nolint:forcetypeassert //on the cases loadValue is used, the value is always of the correct type
	}

	var zero T

	return zero
}

func (rtt *RoundTripTracer) updateStatusBitsAndNotify(flag TraceEvent) {
	previous := rtt.statusBits.Or(uint64(flag))
	updatedBits := previous | uint64(flag)

	if (updatedBits&uint64(rtt.notificationMask)) != 0 && rtt.notificationChan != nil {
		rtt.notificationChan <- rtt.Status()
	}
}
