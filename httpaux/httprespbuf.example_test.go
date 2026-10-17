package httpaux_test

import (
	"fmt"
	"io"
	"net/http/httptest"

	"github.com/angrifel/unapologetic/httpaux"
)

func ExampleBufferResponseBody() {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(200)
	_, _ = recorder.Write([]byte("buffered content"))

	resp := httpaux.BufferResponseBody(recorder.Result())

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Seek back to the start and read the body again.
	_, _ = resp.Body.(io.Seeker).Seek(0, io.SeekStart)
	body, _ = io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// buffered content
	// buffered content
}
