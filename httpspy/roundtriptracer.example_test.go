package httpspy_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/angrifel/unapologetic/httpspy"
)

func ExampleNewRoundTripTracer() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}))
	defer server.Close()

	tracer := httpspy.NewRoundTripTracer()
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req = tracer.TraceRequest(req)
	resp, _ := http.DefaultTransport.RoundTrip(req)

	tracedResp := tracer.TraceResponse(resp)
	body, _ := io.ReadAll(tracedResp.Body)
	tracedResp.Body.Close()

	status := tracer.Status()
	fmt.Println(string(body))
	fmt.Println(status.StatusCode)
	fmt.Println(status.ConnectionTime > 0)
	// Output:
	// hello
	// 200
	// true
}

func ExampleNewRoundTripTracerWithNotification() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ch := make(chan httpspy.RoundTripTracerStatus, 10)
	tracer := httpspy.NewRoundTripTracerWithNotification(httpspy.ResponseHeaders, ch)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req = tracer.TraceRequest(req)
	resp, _ := http.DefaultTransport.RoundTrip(req)
	tracedResp := tracer.TraceResponse(resp)
	tracedResp.Body.Close()

	notification := <-ch
	fmt.Println(notification.StatusCode)
	// Output:
	// 200
}
