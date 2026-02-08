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

	// Body can now be read multiple times via seeking
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// buffered content
}
