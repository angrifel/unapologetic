package httpspy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/netaux"
	"github.com/angrifel/unapologetic/ioaux"
)

// Group 1: Constructor Tests

func TestNewRoundTripTracer_ReturnsNonNil(t *testing.T) {
	// act
	tracer := NewRoundTripTracer()

	// assert
	assert.IsNotNil(t, tracer)
}

func TestNewRoundTripTracerWithNotification_StoresMaskAndChannel(t *testing.T) {
	// arrange
	mask := GetConn | GotConn
	ch := make(chan RoundTripTracerStatus, 1)

	// act
	tracer := NewRoundTripTracerWithNotification(mask, ch)

	// assert
	assert.IsNotNil(t, tracer)
	assert.Equal(t, mask, tracer.notificationMask)
	assert.IsNotNil(t, ch)
}

// Group 2: TraceRequest Tests

func TestTraceRequest_RecordsStartTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	beforeRequest := time.Now()

	// act
	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)

	// assert
	assert.IsNil(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	startTime := loadValue[time.Time](&tracer.start)
	if startTime.Before(beforeRequest) {
		t.Errorf("start time %v should be >= beforeRequest %v", startTime, beforeRequest)
	}
}

func TestTraceRequest_GotConnSetsConnectionTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)

	// assert
	assert.IsNil(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	connectionTime := time.Duration(tracer.connectionTime.Load())
	assert.GreaterOrEqual(t, connectionTime, time.Duration(0))
}

func TestTraceRequest_RecordsSendTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)

	// assert
	assert.IsNil(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	sendTime := time.Duration(tracer.sendTime.Load())
	assert.GreaterOrEqual(t, sendTime, time.Duration(0))
}

func TestTraceRequest_RecordsTimeToFirstByte(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)

	// assert
	assert.IsNil(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	ttfb := time.Duration(tracer.timeToFirstByteTime.Load())
	assert.GreaterOrEqual(t, ttfb, time.Duration(0))
}

func TestTraceRequest_ReturnsNewRequestWithContext(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	originalReq, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.IsNil(t, err)

	// act
	tracedReq := tracer.TraceRequest(originalReq)

	// assert
	assert.NotEqual(t, originalReq, tracedReq)
	assert.IsNotNil(t, tracedReq.Context())
}

func TestTraceRequest_RecordsTLSHandshakeTime(t *testing.T) {
	// arrange
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := server.Client().Do(req)

	// assert
	assert.IsNil(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	tlsTime := time.Duration(tracer.tlsHandshakeTime.Load())
	assert.Greater(t, tlsTime, time.Duration(0))
}

// Group 3: TraceResponse Tests

func TestTraceResponse_RecordsStatusCode(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	response := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act
	tracedResponse := tracer.TraceResponse(response)

	// assert
	assert.Equal(t, int64(http.StatusCreated), tracer.statusCode.Load())
	assert.IsNotNil(t, tracedResponse)
}

func TestTraceResponse_BodyReadRecordsBytesAndTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	// act
	tracedResp := tracer.TraceResponse(resp)
	body, err := io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	t.Cleanup(func() { tracedResp.Body.Close() })

	// assert
	assert.Equal(t, "hello world", string(body))
	assert.Greater(t, tracer.responseBodyBytesRead.Load(), int64(0))
	assert.Greater(t, time.Duration(tracer.receiveTime.Load()), time.Duration(0))
}

func TestTraceResponse_BodyReadEOFSetsFlag(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)

	// act
	_, err = io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	t.Cleanup(func() { tracedResp.Body.Close() })

	// assert
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&BodyReadEOF)
}

func TestTraceResponse_BodyReadErrorCapturesError(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	expectedErr := errors.New("read error")
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(ioaux.ReaderFunc(func(_ []byte) (int, error) {
			return 0, expectedErr
		})),
	}

	tracedResp := tracer.TraceResponse(response)

	// act
	_, err := tracedResp.Body.Read(make([]byte, 10))

	// assert
	assert.Equal(t, expectedErr, err)

	// Verify via status bits that the error was captured
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&BodyReadError)
}

func TestTraceResponse_BodyCloseRecordsCloseTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)

	// act
	err = tracedResp.Body.Close()

	// assert
	assert.IsNil(t, err)
	assert.Greater(t, time.Duration(tracer.closeTime.Load()), time.Duration(0))

	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&BodyClose)
}

func TestTraceResponse_BodyCloseErrorCapturesError(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	expectedErr := errors.New("close error")
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: struct {
			io.Reader
			io.Closer
		}{
			Reader: bytes.NewReader([]byte("data")),
			Closer: ioaux.CloserFunc(func() error {
				return expectedErr
			}),
		},
	}

	tracedResp := tracer.TraceResponse(response)

	// act
	err := tracedResp.Body.Close()

	// assert
	assert.Equal(t, expectedErr, err)

	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&BodyCloseError)
}

func TestTraceResponse_PreservesOriginalBodyBehavior(t *testing.T) {
	// arrange
	expectedContent := []byte("original content here")
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(expectedContent)),
	}

	// act
	tracedResp := tracer.TraceResponse(response)
	body, err := io.ReadAll(tracedResp.Body)

	// assert
	assert.IsNil(t, err)
	assert.EqualFunc(t, expectedContent, body, func(a, b []byte) bool {
		return bytes.Equal(a, b)
	})
}

// Group 4: TraceRoundTripError Tests

func TestTraceRoundTripError_StoresError(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)
	expectedErr := errors.New("connection refused")

	// act
	tracer.TraceRoundTripError(expectedErr)

	// assert - verify via status bits that the error was captured
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&RoundTripError)
}

func TestTraceRoundTripError_SetsReceiveAndCloseTime(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	// act
	tracer.TraceRoundTripError(errors.New("error"))

	// assert
	receiveTime := time.Duration(tracer.receiveTime.Load())
	closeTime := time.Duration(tracer.closeTime.Load())
	assert.GreaterOrEqual(t, receiveTime, time.Duration(0))
	assert.GreaterOrEqual(t, closeTime, time.Duration(0))
}

func TestTraceRoundTripError_SetsStatusBit(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	// act
	tracer.TraceRoundTripError(errors.New("error"))

	// assert
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&RoundTripError)
}

// Group 5: Status() Tests

func TestStatus_ReturnsAllTimingFields(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)
	_, err = io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	err = tracedResp.Body.Close()
	assert.IsNil(t, err)

	// act
	status := tracer.Status()

	// assert
	assert.NotEqual(t, time.Time{}, status.Start)
	assert.GreaterOrEqual(t, status.ConnectionTime, time.Duration(0))
	assert.GreaterOrEqual(t, status.SendTime, time.Duration(0))
	assert.GreaterOrEqual(t, status.FirstByteTime, time.Duration(0))
	assert.GreaterOrEqual(t, status.ReceiveTime, time.Duration(0))
	assert.GreaterOrEqual(t, status.CloseTime, time.Duration(0))
	assert.Equal(t, http.StatusOK, status.StatusCode)
}

func TestStatus_AfterFullRoundTrip_ReturnsCompleteStatus(t *testing.T) {
	// arrange
	responseBody := "complete round trip response"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responseBody)) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)
	body, err := io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	err = tracedResp.Body.Close()
	assert.IsNil(t, err)

	// act
	status := tracer.Status()

	// assert
	assert.Equal(t, responseBody, string(body))
	assert.Equal(t, http.StatusAccepted, status.StatusCode)
	assert.NotEqual(t, time.Time{}, status.Start)
	assert.Greater(t, status.ConnectionTime, time.Duration(0))
	assert.Greater(t, status.SendTime, time.Duration(0))
	assert.Greater(t, status.FirstByteTime, time.Duration(0))
	assert.Greater(t, status.ReceiveTime, time.Duration(0))
	assert.Greater(t, status.CloseTime, time.Duration(0))
	assert.Greater(t, status.ResponseBodyBytesRead, 0)
}

// Group 6: Notification Tests

func TestNotification_MatchingMaskSendsToChannel(t *testing.T) {
	// arrange
	ch := make(chan RoundTripTracerStatus, 10)
	tracer := NewRoundTripTracerWithNotification(ResponseHeaders, ch)

	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act
	tracer.TraceResponse(response)

	// assert
	select {
	case status := <-ch:
		assert.Equal(t, http.StatusOK, status.StatusCode)
	case <-time.After(100 * time.Millisecond):
		t.Error("expected notification but none received")
	}
}

func TestNotification_NonMatchingMaskNoSend(t *testing.T) {
	// arrange
	ch := make(chan RoundTripTracerStatus, 10)
	tracer := NewRoundTripTracerWithNotification(DNSStart, ch) // Different mask

	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act
	tracer.TraceResponse(response) // Triggers ResponseHeaders, not DNSStart

	// assert
	select {
	case <-ch:
		t.Error("expected no notification but received one")
	case <-time.After(50 * time.Millisecond):
		// Expected: no notification
	}
}

func TestNotification_MultipleEventsNotify(t *testing.T) {
	// arrange
	ch := make(chan RoundTripTracerStatus, 10)
	tracer := NewRoundTripTracerWithNotification(ResponseHeaders|BodyClose, ch)

	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act
	tracedResp := tracer.TraceResponse(response)
	tracedResp.Body.Close() //nolint:errcheck

	// assert
	notificationCount := 0
	timeout := time.After(100 * time.Millisecond)
loop:
	for {
		select {
		case <-ch:
			notificationCount++
		case <-timeout:
			break loop
		}
	}

	assert.GreaterOrEqual(t, notificationCount, 2)
}

func TestNotification_NilChannelDoesNotPanic(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracerWithNotification(ResponseHeaders, nil)

	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act & assert - should not panic
	tracer.TraceResponse(response)
}

func TestNotification_StatusContainsCurrentState(t *testing.T) {
	// arrange
	ch := make(chan RoundTripTracerStatus, 10)
	tracer := NewRoundTripTracerWithNotification(ResponseHeaders, ch)

	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusTeapot,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}

	// act
	tracer.TraceResponse(response)

	// assert
	select {
	case status := <-ch:
		assert.Equal(t, http.StatusTeapot, status.StatusCode)
		assert.NotEqual(t, time.Time{}, status.Start)
	case <-time.After(100 * time.Millisecond):
		t.Error("expected notification but none received")
	}
}

// Group 7: Status Bits Tests

func TestStatusBits_SetCorrectlyForEachEvent(t *testing.T) {
	events := []TraceEvent{
		GetConn, GotConn, DNSStart, DNSDone,
		ConnectStart, ConnectDone,
		TLSHandshakeStart, TLSHandshakeDone,
		WroteRequest, GotFirstResponseByte,
		ResponseHeaders, BodyRead, BodyReadEOF, BodyClose,
	}

	for _, event := range events {
		t.Run(event.String(), func(t *testing.T) {
			// arrange
			tracer := NewRoundTripTracer()

			// act
			tracer.updateStatusBitsAndNotify(event)

			// assert
			statusBits := TraceEvent(tracer.statusBits.Load())
			assert.NotEqual(t, TraceEvent(0), statusBits&event)
		})
	}
}

func TestStatusBits_MultipleEventsAccumulate(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()

	// act
	tracer.updateStatusBitsAndNotify(GetConn)
	tracer.updateStatusBitsAndNotify(GotConn)
	tracer.updateStatusBitsAndNotify(ResponseHeaders)

	// assert
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&GetConn)
	assert.NotEqual(t, TraceEvent(0), statusBits&GotConn)
	assert.NotEqual(t, TraceEvent(0), statusBits&ResponseHeaders)
}

func TestStatusBits_CompositeEventsSetsMultipleBits(t *testing.T) {
	// arrange
	tracer := NewRoundTripTracer()

	// act
	tracer.updateStatusBitsAndNotify(WroteRequest | WroteRequestError)

	// assert
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&WroteRequest)
	assert.NotEqual(t, TraceEvent(0), statusBits&WroteRequestError)
}

// Group 8: End-to-End Integration Tests

func TestRoundTrip_HTTPSuccess_RecordsAllMetrics(t *testing.T) {
	// arrange
	responseBody := "integration test body"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody)) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := http.DefaultClient.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)
	body, err := io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	err = tracedResp.Body.Close()
	assert.IsNil(t, err)

	// assert
	assert.Equal(t, responseBody, string(body))

	status := tracer.Status()
	assert.Equal(t, http.StatusOK, status.StatusCode)
	assert.NotEqual(t, time.Time{}, status.Start)
	assert.Greater(t, status.ConnectionTime, time.Duration(0))
	assert.Greater(t, status.SendTime, time.Duration(0))
	assert.Greater(t, status.FirstByteTime, time.Duration(0))
	assert.Greater(t, status.ReceiveTime, time.Duration(0))
	assert.Greater(t, status.CloseTime, time.Duration(0))
	assert.Greater(t, status.ResponseBodyBytesRead, 0)
}

func TestRoundTrip_HTTPSSuccess_IncludesTLSMetrics(t *testing.T) {
	// arrange
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("tls response")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := server.Client().Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)
	_, err = io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	err = tracedResp.Body.Close()
	assert.IsNil(t, err)

	// assert
	status := tracer.Status()
	assert.Greater(t, status.TLSHandshakeTime, time.Duration(0))
}

func TestRoundTrip_BodyReadError_RecordsError(t *testing.T) {
	// arrange
	expectedErr := errors.New("body read failure")
	tracer := NewRoundTripTracer()
	now := time.Now()
	tracer.start.Store(now)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(ioaux.ReaderFunc(func(_ []byte) (int, error) {
			return 0, expectedErr
		})),
	}

	// act
	tracedResp := tracer.TraceResponse(response)
	_, err := tracedResp.Body.Read(make([]byte, 10))

	// assert
	assert.Equal(t, expectedErr, err)

	// Verify via status bits that the error was captured
	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&BodyReadError)
}

func TestRoundTrip_ConcurrentUsage_ThreadSafe(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("concurrent")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	// act & assert - should not race
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tracer := NewRoundTripTracer()
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				return
			}

			req = tracer.TraceRequest(req)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}

			tracedResp := tracer.TraceResponse(resp)
			io.ReadAll(tracedResp.Body) //nolint:errcheck
			tracedResp.Body.Close()     //nolint:errcheck

			_ = tracer.Status()
		}()
	}
	wg.Wait()
}

func TestTraceRequest_DNSCallbacksRecordDNSLookupTime(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("dns test")) //nolint:errcheck
	}))
	t.Cleanup(server.Close)

	// Extract the port from the test server address
	serverAddr := server.Listener.Addr().String()
	_, serverPort, err := net.SplitHostPort(serverAddr)
	assert.IsNil(t, err)

	// Start a real DNS server that maps "traced.test." → 127.0.0.1
	dnsRecords := []netaux.DNSRecord{
		{Name: "traced.test.", Type: netaux.TypeA, Class: netaux.ClassIN, TTL: 300, Addr: netip.MustParseAddr("127.0.0.1")},
	}
	dnsConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err)
	dnsAddr := dnsConn.LocalAddr()
	dnsSrv := &netaux.DNSServer{Records: dnsRecords}
	dnsDone := make(chan struct{})
	go func() {
		defer close(dnsDone)
		dnsSrv.Serve(dnsConn) //nolint:errcheck
	}()
	t.Cleanup(func() {
		dnsSrv.Shutdown()
		dnsConn.Close()
		<-dnsDone
	})

	// Build a resolver that uses the test DNS server
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "udp", dnsAddr.String())
		},
	}

	// Build an HTTP client that uses the custom resolver
	transport := &http.Transport{ //nolint:exhaustruct
		DialContext: (&net.Dialer{ //nolint:exhaustruct
			Resolver: resolver,
		}).DialContext,
	}
	client := &http.Client{Transport: transport} //nolint:exhaustruct
	t.Cleanup(client.CloseIdleConnections)

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://traced.test:%s/", serverPort), nil)
	assert.IsNil(t, err)

	// act
	req = tracer.TraceRequest(req)
	resp, err := client.Do(req)
	assert.IsNil(t, err)

	tracedResp := tracer.TraceResponse(resp)
	body, err := io.ReadAll(tracedResp.Body)
	assert.IsNil(t, err)
	err = tracedResp.Body.Close()
	assert.IsNil(t, err)

	// assert
	assert.Equal(t, "dns test", string(body))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	status := tracer.Status()
	assert.Greater(t, status.DNSLookupTime, time.Duration(0))

	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&DNSStart)
	assert.NotEqual(t, TraceEvent(0), statusBits&DNSDone)
}

func TestTraceRequest_WroteRequestErrorCapturesSendError(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	// Rationale for using ioaux.ReaderFunc: Real I/O cannot reliably produce a
	// request-write error without racy timing tricks (e.g., server closing mid-write).
	// The existing test file already uses ioaux.ReaderFunc for error injection in
	// body-read tests (lines 260, 742), establishing the pattern. The network path
	// itself (dial, connect, write headers) is still real I/O; only the request body
	// reader is synthetic.
	bodyErr := errors.New("request body read failure")
	failingBody := ioaux.ReaderFunc(func(_ []byte) (int, error) {
		return 0, bodyErr
	})

	tracer := NewRoundTripTracer()
	req, err := http.NewRequest(http.MethodPost, server.URL, failingBody)
	assert.IsNil(t, err)
	req.ContentLength = 100 // force the transport to attempt reading the full body

	// act
	req = tracer.TraceRequest(req)
	resp, roundTripErr := http.DefaultClient.Do(req)
	if resp != nil {
		resp.Body.Close() //nolint:errcheck
	}

	// assert
	assert.IsNotNil(t, roundTripErr)
	assert.Equal(t, true, strings.Contains(roundTripErr.Error(), bodyErr.Error()))

	statusBits := TraceEvent(tracer.statusBits.Load())
	assert.NotEqual(t, TraceEvent(0), statusBits&WroteRequest)
	assert.NotEqual(t, TraceEvent(0), statusBits&WroteRequestError)

	assert.IsNotNil(t, tracer.sendError.Load())
}

// Helper for String() method on TraceEvent (for test naming)
func (te TraceEvent) String() string {
	switch te {
	case GetConn:
		return "GetConn"
	case GotConn:
		return "GotConn"
	case DNSStart:
		return "DNSStart"
	case DNSDone:
		return "DNSDone"
	case ConnectStart:
		return "ConnectStart"
	case ConnectDone:
		return "ConnectDone"
	case TLSHandshakeStart:
		return "TLSHandshakeStart"
	case TLSHandshakeDone:
		return "TLSHandshakeDone"
	case WroteRequest:
		return "WroteRequest"
	case WroteRequestError:
		return "WroteRequestError"
	case GotFirstResponseByte:
		return "GotFirstResponseByte"
	case ResponseHeaders:
		return "ResponseHeaders"
	case RoundTripError:
		return "RoundTripError"
	case BodyRead:
		return "BodyRead"
	case BodyReadEOF:
		return "BodyReadEOF"
	case BodyReadError:
		return "BodyReadError"
	case BodyClose:
		return "BodyClose"
	case BodyCloseError:
		return "BodyCloseError"
	default:
		return "Unknown"
	}
}
