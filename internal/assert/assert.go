// Package assert provides type-safe assertion functions for testing.
package assert

import (
	"cmp"
	"fmt"
	"reflect"
)

type testingT = interface {
	Errorf(format string, args ...interface{})
	Helper()
}

// IsNil asserts that the value is nil.
func IsNil(t testingT, value any, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if nil != value && !reflect.ValueOf(value).IsNil() {
		t.Errorf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, nil, value)
	}
}

// IsNotNil asserts that the value is not nil.
func IsNotNil(t testingT, value any, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if nil == value || reflect.ValueOf(value).IsNil() {
		t.Errorf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, "not nil", value)
	}
}

// Equal asserts that two values are equal using the == operator.
func Equal[T comparable](t testingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if expected != actual {
		t.Errorf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, expected, actual)
	}
}

// NotEqual asserts that two values are not equal using the != operator.
func NotEqual[T comparable](t testingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if expected == actual {
		t.Errorf("%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------", prefix, expected)
	}
}

// EqualFunc asserts that two values are equal using the provided equality function.
func EqualFunc[T any](t testingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if !equalFunc(expected, actual) {
		t.Errorf("%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, expected, actual)
	}
}

// NotEqualFunc asserts that two values are not equal using the provided equality function.
func NotEqualFunc[T any](t testingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if equalFunc(expected, actual) {
		t.Errorf("%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------", prefix, expected)
	}
}

// Less asserts that the actual value is less than the expected value.
func Less[T cmp.Ordered](t testingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual >= maxValue {
		t.Errorf("%s\nexpected:----------------------------\n%v to be less than %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, maxValue, actual)
	}
}

// LessOrEqual asserts that the actual value is less than or equal to the expected value.
func LessOrEqual[T cmp.Ordered](t testingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual > maxValue {
		t.Errorf("%s\nexpected:----------------------------\n%v to be less than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, maxValue, actual)
	}
}

// Greater asserts that the actual value is greater than the expected value.
func Greater[T cmp.Ordered](t testingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual <= minValue {
		t.Errorf("%s\nexpected:----------------------------\n%v to be greater than %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, minValue, actual)
	}
}

// GreaterOrEqual asserts that the actual value is greater than or equal to the expected value.
func GreaterOrEqual[T cmp.Ordered](t testingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	prefix := ""
	if len(messageAndArgs) > 0 {
		prefix = fmt.Sprintf(fmt.Sprint(messageAndArgs[0]), messageAndArgs[1:]...)
	}

	if actual < minValue {
		t.Errorf("%s\nexpected:----------------------------\n%v to be greater than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------", prefix, actual, minValue, actual)
	}
}
