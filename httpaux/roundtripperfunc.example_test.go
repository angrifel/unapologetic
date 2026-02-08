package httpaux_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/angrifel/unapologetic/httpaux"
)

func ExampleRoundTripperFunc() {
	rt := httpaux.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("mock response")),
		}, nil
	})

	client := &http.Client{Transport: rt}
	resp, _ := client.Get("http://example.com")
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	// Output:
	// 200
	// mock response
}
