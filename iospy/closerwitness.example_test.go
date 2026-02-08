package iospy_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/angrifel/unapologetic/iospy"
)

func ExampleWitnessCloser() {
	closer := io.NopCloser(strings.NewReader("data"))
	witnessed := iospy.WitnessCloser(closer)

	witnessed.Close() //nolint:errcheck

	calls := witnessed.(iospy.CloserWitness).ObservedCloseCalls()
	fmt.Println(len(calls))
	fmt.Println(calls[0].ResultErr)
	// Output:
	// 1
	// <nil>
}
