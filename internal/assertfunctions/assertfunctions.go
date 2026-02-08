// Package assertfunctions provides type-safe assertion functions for testing.
package assertfunctions

import (
	"cmp"
	"fmt"
	"reflect"
)

type TestingTCommon = interface {
	Logf(format string, args ...interface{})
	Helper()
}

type TestingT struct {
	TestingTCommon
	Fail func()
}

// IsNil asserts that the value is nil.
func IsNil(t TestingT, value any, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if nil != value && !reflect.ValueOf(value).IsNil() {
		t.Logf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, nil, value)
		t.Fail()
	}
}

// IsNotNil asserts that the value is not nil.
func IsNotNil(t TestingT, value any, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if nil == value || reflect.ValueOf(value).IsNil() {
		t.Logf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, "not nil", value)
		t.Fail()
	}
}

// Equal asserts that two values are equal using the == operator.
func Equal[T comparable](t TestingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if expected != actual {
		t.Logf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, expected, actual)
		t.Fail()
	}
}

// NotEqual asserts that two values are not equal using the != operator.
func NotEqual[T comparable](t TestingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if expected == actual {
		t.Logf("%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------", prefix, expected)
		t.Fail()
	}
}

// EqualFunc asserts that two values are equal using the provided equality function.
func EqualFunc[T any](t TestingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if !equalFunc(expected, actual) {
		t.Logf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, expected, actual)
		t.Fail()
	}
}

// NotEqualFunc asserts that two values are not equal using the provided equality function.
func NotEqualFunc[T any](t TestingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if equalFunc(expected, actual) {
		t.Logf("%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------", prefix, expected)
		t.Fail()
	}
}

// Less asserts that the actual value is less than the expected value.
func Less[T cmp.Ordered](t TestingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual >= maxValue {
		t.Logf("%s\nexpected:----------------------------\n%v to be less than %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, maxValue, actual)
		t.Fail()
	}
}

// LessOrEqual asserts that the actual value is less than or equal to the expected value.
func LessOrEqual[T cmp.Ordered](t TestingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual > maxValue {
		t.Logf("%s\nexpected:----------------------------\n%v to be less than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, maxValue, actual)
		t.Fail()
	}
}

// Greater asserts that the actual value is greater than the expected value.
func Greater[T cmp.Ordered](t TestingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual <= minValue {
		t.Logf("%s\nexpected:----------------------------\n%v to be greater than %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, minValue, actual)
		t.Fail()
	}
}

// GreaterOrEqual asserts that the actual value is greater than or equal to the expected value.
func GreaterOrEqual[T cmp.Ordered](t TestingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual < minValue {
		t.Logf("%s\nexpected:----------------------------\n%v to be greater than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, minValue, actual)
		t.Fail()
	}
}
