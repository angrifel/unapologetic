// Package assert provides type-safe assertion functions for testing.
package assert

import (
	"cmp"

	"github.com/angrifel/unapologetic/internal/assertfunctions"
)

// TestingT is an interface that defines a minimal set of methods for logging, marking failures, and helper annotations.
type TestingT = interface {
	Fail()
	Helper()
	Logf(format string, args ...interface{})
}

// IsNil asserts that the value is nil.
func IsNil(t TestingT, value any, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.IsNil(aft, value, messageAndArgs...)
}

// IsNotNil asserts that the value is not nil.
func IsNotNil(t TestingT, value any, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.IsNotNil(aft, value, messageAndArgs...)
}

// Equal asserts that two values are equal using the == operator.
func Equal[T comparable](t TestingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.Equal(aft, expected, actual, messageAndArgs...)
}

// NotEqual asserts that two values are not equal using the != operator.
func NotEqual[T comparable](t TestingT, expected, actual T, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.NotEqual(aft, expected, actual, messageAndArgs...)
}

// EqualFunc asserts that two values are equal using the provided equality function.
func EqualFunc[T any](t TestingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.EqualFunc(aft, expected, actual, equalFunc, messageAndArgs...)
}

// NotEqualFunc asserts that two values are not equal using the provided equality function.
func NotEqualFunc[T any](t TestingT, expected, actual T, equalFunc func(a, b T) bool, messageAndArgs ...interface{}) {
	t.Helper()
	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.NotEqualFunc(aft, expected, actual, equalFunc, messageAndArgs...)
}

// Less asserts that the actual value is less than maxValue.
func Less[T cmp.Ordered](t TestingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.Less(aft, actual, maxValue, messageAndArgs...)
}

// LessOrEqual asserts that the actual value is less than or equal to maxValue.
func LessOrEqual[T cmp.Ordered](t TestingT, actual, maxValue T, messageAndArgs ...interface{}) {
	t.Helper()

	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.LessOrEqual(aft, actual, maxValue, messageAndArgs...)
}

// Greater asserts that the actual value is greater than minValue.
func Greater[T cmp.Ordered](t TestingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.Greater(aft, actual, minValue, messageAndArgs...)
}

// GreaterOrEqual asserts that the actual value is greater than or equal to minValue.
func GreaterOrEqual[T cmp.Ordered](t TestingT, actual, minValue T, messageAndArgs ...interface{}) {
	t.Helper()

	aft := assertfunctions.TestingT{TestingTCommon: t, Fail: t.Fail}
	assertfunctions.GreaterOrEqual(aft, actual, minValue, messageAndArgs...)
}
