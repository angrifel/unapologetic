package assertfunctions

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

// testingTSpy wraps *testing.T to capture Logf and Fail calls made by assertion functions.
// It shadows Logf to record log output without writing to the test log, while Helper()
// passes through to the embedded *testing.T.
type testingTSpy struct {
	*testing.T
	logCalled  bool
	logMessage string
	failCalled bool
}

func (s *testingTSpy) Logf(format string, args ...interface{}) {
	s.logCalled = true
	// We intentionally do NOT call s.T.Logf to avoid polluting test output
	// with expected assertion failures.
	s.logMessage = fmt.Sprintf(format, args...)
}

func newTestingT(t *testing.T, useFailNow bool) (*testingTSpy, TestingT) {
	t.Helper()
	spy := &testingTSpy{T: t}
	fail := func() { spy.failCalled = true }
	if useFailNow {
		fail = func() {
			spy.failCalled = true
			runtime.Goexit()
		}
	}
	return spy, TestingT{TestingTCommon: spy, Fail: fail}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"nil interface", nil, nil, false, false},
		{"nil pointer", (*int)(nil), nil, false, false},
		{"non-nil pointer", new(int), nil, true, false},
		{"non-nil pointer with message", new(int), []interface{}{"val %d", 1}, true, false},
		// One nil and one non-nil value per nillable kind, plus non-nillable kinds.
		// reflect.Interface is not represented: a value passed as `any` is always
		// reported by reflect.ValueOf under its concrete kind.
		{"nil chan", (chan int)(nil), nil, false, false},
		{"non-nil chan", make(chan int), nil, true, false},
		{"nil func", (func())(nil), nil, false, false},
		{"non-nil func", func() {}, nil, true, false},
		{"nil map", (map[int]int)(nil), nil, false, false},
		{"non-nil map", map[int]int{}, nil, true, false},
		{"nil slice", ([]int)(nil), nil, false, false},
		{"non-nil slice", []int{}, nil, true, false},
		{"nil unsafe pointer", unsafe.Pointer(nil), nil, false, false},
		{"non-nil unsafe pointer", unsafe.Pointer(new(int)), nil, true, false},
		{"int", 0, nil, true, false},
		{"string", "", nil, true, false},
		{"struct", struct{}{}, nil, true, false},
		{"array", [1]int{}, nil, true, false},
		{"nil interface/failnow", nil, nil, false, true},
		{"nil pointer/failnow", (*int)(nil), nil, false, true},
		{"non-nil pointer/failnow", new(int), nil, true, true},
		{"int/failnow", 0, nil, true, true},
		{"non-nil pointer with message/failnow", new(int), []interface{}{"val %d", 1}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				IsNil(tt, tc.value, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestIsNotNil(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"non-nil pointer", new(int), nil, false, false},
		{"nil interface", nil, nil, true, false},
		{"nil pointer", (*int)(nil), nil, true, false},
		{"nil interface with message", nil, []interface{}{"should not be nil"}, true, false},
		// One nil and one non-nil value per nillable kind, plus non-nillable kinds.
		// reflect.Interface is not represented: a value passed as `any` is always
		// reported by reflect.ValueOf under its concrete kind.
		{"nil chan", (chan int)(nil), nil, true, false},
		{"non-nil chan", make(chan int), nil, false, false},
		{"nil func", (func())(nil), nil, true, false},
		{"non-nil func", func() {}, nil, false, false},
		{"nil map", (map[int]int)(nil), nil, true, false},
		{"non-nil map", map[int]int{}, nil, false, false},
		{"nil slice", ([]int)(nil), nil, true, false},
		{"non-nil slice", []int{}, nil, false, false},
		{"nil unsafe pointer", unsafe.Pointer(nil), nil, true, false},
		{"non-nil unsafe pointer", unsafe.Pointer(new(int)), nil, false, false},
		{"int", 0, nil, false, false},
		{"string", "", nil, false, false},
		{"struct", struct{}{}, nil, false, false},
		{"array", [1]int{}, nil, false, false},
		{"non-nil pointer/failnow", new(int), nil, false, true},
		{"int/failnow", 0, nil, false, true},
		{"nil interface/failnow", nil, nil, true, true},
		{"nil pointer/failnow", (*int)(nil), nil, true, true},
		{"nil interface with message/failnow", nil, []interface{}{"should not be nil"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				IsNotNil(tt, tc.value, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestMessageAndArgs(t *testing.T) {
	equal := func(a, b int) bool { return a == b }

	// Every case triggers a failure so the prefix is logged.
	tests := []struct {
		name       string
		assertion  func(tt TestingT, msg ...interface{})
		msg        []interface{}
		wantPrefix string
	}{
		{"IsNil/no message", func(tt TestingT, msg ...interface{}) { IsNil(tt, 0, msg...) }, nil, "\n"},
		{"IsNil/format with args", func(tt TestingT, msg ...interface{}) { IsNil(tt, 0, msg...) }, []interface{}{"val %d", 1}, "val 1\n"},
		{"IsNotNil/message only", func(tt TestingT, msg ...interface{}) { IsNotNil(tt, nil, msg...) }, []interface{}{"plain"}, "plain\n"},
		{"IsNotNil/non-string message", func(tt TestingT, msg ...interface{}) { IsNotNil(tt, nil, msg...) }, []interface{}{42}, "42\n"},
		{"Equal", func(tt TestingT, msg ...interface{}) { Equal(tt, 1, 2, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"NotEqual", func(tt TestingT, msg ...interface{}) { NotEqual(tt, 1, 1, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"EqualFunc", func(tt TestingT, msg ...interface{}) { EqualFunc(tt, 1, 2, equal, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"NotEqualFunc", func(tt TestingT, msg ...interface{}) { NotEqualFunc(tt, 1, 1, equal, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"Less", func(tt TestingT, msg ...interface{}) { Less(tt, 2, 1, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"LessOrEqual", func(tt TestingT, msg ...interface{}) { LessOrEqual(tt, 2, 1, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"Greater", func(tt TestingT, msg ...interface{}) { Greater(tt, 1, 2, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
		{"GreaterOrEqual", func(tt TestingT, msg ...interface{}) { GreaterOrEqual(tt, 1, 2, msg...) }, []interface{}{"item %d", 0}, "item 0\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spy, tt := newTestingT(t, false)
			tc.assertion(tt, tc.msg...)

			if !strings.HasPrefix(spy.logMessage, tc.wantPrefix) {
				t.Errorf("logMessage = %q, want prefix %q", spy.logMessage, tc.wantPrefix)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name       string
		expected   int
		actual     int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"equal values", 42, 42, nil, false, false},
		{"unequal values", 1, 2, nil, true, false},
		{"unequal with message", 1, 2, []interface{}{"item %d", 0}, true, false},
		{"equal values/failnow", 42, 42, nil, false, true},
		{"unequal values/failnow", 1, 2, nil, true, true},
		{"unequal with message/failnow", 1, 2, []interface{}{"item %d", 0}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				Equal(tt, tc.expected, tc.actual, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name       string
		expected   int
		actual     int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"different values", 1, 2, nil, false, false},
		{"same values", 42, 42, nil, true, false},
		{"same values with message", 42, 42, []interface{}{"should differ"}, true, false},
		{"different values/failnow", 1, 2, nil, false, true},
		{"same values/failnow", 42, 42, nil, true, true},
		{"same values with message/failnow", 42, 42, []interface{}{"should differ"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				NotEqual(tt, tc.expected, tc.actual, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestEqualFunc(t *testing.T) {
	alwaysEqual := func(_, _ string) bool { return true }
	neverEqual := func(_, _ string) bool { return false }

	tests := []struct {
		name       string
		expected   string
		actual     string
		equalFunc  func(string, string) bool
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"func returns true", "a", "a", alwaysEqual, nil, false, false},
		{"func returns false", "a", "b", neverEqual, nil, true, false},
		{"func returns false with message", "a", "b", neverEqual, []interface{}{"compare %s", "failed"}, true, false},
		{"func returns true/failnow", "a", "a", alwaysEqual, nil, false, true},
		{"func returns false/failnow", "a", "b", neverEqual, nil, true, true},
		{"func returns false with message/failnow", "a", "b", neverEqual, []interface{}{"compare %s", "failed"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				EqualFunc(tt, tc.expected, tc.actual, tc.equalFunc, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestNotEqualFunc(t *testing.T) {
	neverEqual := func(_, _ string) bool { return false }
	alwaysEqual := func(_, _ string) bool { return true }

	tests := []struct {
		name       string
		expected   string
		actual     string
		equalFunc  func(string, string) bool
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"func returns false (values differ)", "a", "b", neverEqual, nil, false, false},
		{"func returns true (values same)", "a", "a", alwaysEqual, nil, true, false},
		{"func returns true with message", "a", "a", alwaysEqual, []interface{}{"should differ"}, true, false},
		{"func returns false (values differ)/failnow", "a", "b", neverEqual, nil, false, true},
		{"func returns true (values same)/failnow", "a", "a", alwaysEqual, nil, true, true},
		{"func returns true with message/failnow", "a", "a", alwaysEqual, []interface{}{"should differ"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				NotEqualFunc(tt, tc.expected, tc.actual, tc.equalFunc, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		name       string
		actual     int
		maxValue   int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"less than", 3, 5, nil, false, false},
		{"equal to", 5, 5, nil, true, false},
		{"greater than", 7, 5, nil, true, false},
		{"greater with message", 7, 5, []interface{}{"value %d", 7}, true, false},
		{"less than/failnow", 3, 5, nil, false, true},
		{"equal to/failnow", 5, 5, nil, true, true},
		{"greater than/failnow", 7, 5, nil, true, true},
		{"greater with message/failnow", 7, 5, []interface{}{"value %d", 7}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				Less(tt, tc.actual, tc.maxValue, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestLessOrEqual(t *testing.T) {
	tests := []struct {
		name       string
		actual     int
		maxValue   int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"less than", 3, 5, nil, false, false},
		{"equal to", 5, 5, nil, false, false},
		{"greater than", 6, 5, nil, true, false},
		{"greater with message", 6, 5, []interface{}{"too big"}, true, false},
		{"less than/failnow", 3, 5, nil, false, true},
		{"equal to/failnow", 5, 5, nil, false, true},
		{"greater than/failnow", 6, 5, nil, true, true},
		{"greater with message/failnow", 6, 5, []interface{}{"too big"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				LessOrEqual(tt, tc.actual, tc.maxValue, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	tests := []struct {
		name       string
		actual     int
		minValue   int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"greater than", 5, 3, nil, false, false},
		{"equal to", 5, 5, nil, true, false},
		{"less than", 3, 5, nil, true, false},
		{"less with message", 3, 5, []interface{}{"too small"}, true, false},
		{"greater than/failnow", 5, 3, nil, false, true},
		{"equal to/failnow", 5, 5, nil, true, true},
		{"less than/failnow", 3, 5, nil, true, true},
		{"less with message/failnow", 3, 5, []interface{}{"too small"}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				Greater(tt, tc.actual, tc.minValue, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}

func TestGreaterOrEqual(t *testing.T) {
	tests := []struct {
		name       string
		actual     int
		minValue   int
		msg        []interface{}
		wantFail   bool
		useFailNow bool
	}{
		{"greater than", 5, 3, nil, false, false},
		{"equal to", 5, 5, nil, false, false},
		{"less than", 4, 5, nil, true, false},
		{"less with message", 4, 5, []interface{}{"val %d too small", 4}, true, false},
		{"greater than/failnow", 5, 3, nil, false, true},
		{"equal to/failnow", 5, 5, nil, false, true},
		{"less than/failnow", 4, 5, nil, true, true},
		{"less with message/failnow", 4, 5, []interface{}{"val %d too small", 4}, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantReachedAfter := !(tc.wantFail && tc.useFailNow)
			spy, tt := newTestingT(t, tc.useFailNow)
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				GreaterOrEqual(tt, tc.actual, tc.minValue, tc.msg...)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if reachedAfter != wantReachedAfter {
				t.Errorf("reachedAfter = %v, want %v", reachedAfter, wantReachedAfter)
			}
		})
	}
}
