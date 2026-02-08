package iospy_test

import (
	"fmt"
	"strings"

	"github.com/angrifel/unapologetic/iospy"
)

func ExampleWitnessReader() {
	reader := strings.NewReader("hello")
	witnessed := iospy.WitnessReader(reader)

	buf := make([]byte, 5)
	witnessed.Read(buf) //nolint:errcheck

	calls := witnessed.(iospy.ReaderWitness).ObservedReadCalls()
	fmt.Println(len(calls))
	fmt.Println(calls[0].ResultN)
	// Output:
	// 1
	// 5
}
