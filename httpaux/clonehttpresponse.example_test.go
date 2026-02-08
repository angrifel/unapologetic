package httpaux_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/angrifel/unapologetic/httpaux"
)

func ExampleCloneHTTPResponseWithBody() {
	original := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       io.NopCloser(strings.NewReader("original body")),
	}

	newBody := io.NopCloser(strings.NewReader("replaced body"))
	cloned := httpaux.CloneHTTPResponseWithBody(original, newBody)

	body, _ := io.ReadAll(cloned.Body)
	fmt.Println(cloned.StatusCode)
	fmt.Println(string(body))
	// Output:
	// 200
	// replaced body
}
