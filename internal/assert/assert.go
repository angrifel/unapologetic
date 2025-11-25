// Package assert provides type-safe assertion functions for testing.
package assert

import "cmp"

type testingT = interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
}

// Equal asserts that two values are equal using the == operator.
func Equal[T comparable](t testingT, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// NotEqual asserts that two values are not equal using the != operator.
func NotEqual[T comparable](t testingT, expected, actual T) {
	t.Helper()

	if expected == actual {
		t.Errorf("Expected values to be different, but both are %v", expected)
	}
}

// Less asserts that the actual value is less than the expected value.
func Less[T cmp.Ordered](t testingT, actual, maxValue T) {
	t.Helper()

	if actual >= maxValue {
		t.Errorf("Expected %v to be less than %v", actual, maxValue)
	}
}

// LessOrEqual asserts that the actual value is less than or equal to the expected value.
func LessOrEqual[T cmp.Ordered](t testingT, actual, maxValue T) {
	t.Helper()

	if actual > maxValue {
		t.Errorf("Expected %v to be less than or equal to %v", actual, maxValue)
	}
}

// Greater asserts that the actual value is greater than the expected value.
func Greater[T cmp.Ordered](t testingT, actual, minValue T) {
	t.Helper()

	if actual <= minValue {
		t.Errorf("Expected %v to be greater than %v", actual, minValue)
	}
}

// GreaterOrEqual asserts that the actual value is greater than or equal to the expected value.
func GreaterOrEqual[T cmp.Ordered](t testingT, actual, minValue T) {
	t.Helper()

	if actual < minValue {
		t.Errorf("Expected %v to be greater than or equal to %v", actual, minValue)
	}
}
