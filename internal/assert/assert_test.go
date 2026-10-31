package assert

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// testingTSpy wraps *testing.T to capture Fail, FailNow and Logf calls made by assertion functions.
// It shadows those methods so expected assertion failures neither fail nor pollute the real test,
// while Helper() passes through to the embedded *testing.T. FailNow mirrors testing.T by exiting
// the calling goroutine, so assertions must be run on a separate goroutine.
type testingTSpy struct {
	*testing.T
	failCalled    bool
	failNowCalled bool
	logCalled     bool
	logMessage    string
}

func (s *testingTSpy) Fail() { s.failCalled = true }

func (s *testingTSpy) FailNow() {
	s.failNowCalled = true
	runtime.Goexit()
}

func (s *testingTSpy) Logf(format string, args ...interface{}) {
	s.logCalled = true
	s.logMessage = fmt.Sprintf(format, args...)
}

func TestAssertions_FailWithoutStopping(t *testing.T) {
	equal := func(a, b int) bool { return a == b }

	tests := []struct {
		name      string
		assertion func(t TestingT, msg ...interface{})
		wantFail  bool
	}{
		{"IsNil/pass", func(t TestingT, msg ...interface{}) { IsNil(t, nil, msg...) }, false},
		{"IsNil/fail", func(t TestingT, msg ...interface{}) { IsNil(t, new(int), msg...) }, true},
		{"IsNotNil/pass", func(t TestingT, msg ...interface{}) { IsNotNil(t, new(int), msg...) }, false},
		{"IsNotNil/fail", func(t TestingT, msg ...interface{}) { IsNotNil(t, nil, msg...) }, true},
		{"Equal/pass", func(t TestingT, msg ...interface{}) { Equal(t, 1, 1, msg...) }, false},
		{"Equal/fail", func(t TestingT, msg ...interface{}) { Equal(t, 1, 2, msg...) }, true},
		{"NotEqual/pass", func(t TestingT, msg ...interface{}) { NotEqual(t, 1, 2, msg...) }, false},
		{"NotEqual/fail", func(t TestingT, msg ...interface{}) { NotEqual(t, 1, 1, msg...) }, true},
		{"EqualFunc/pass", func(t TestingT, msg ...interface{}) { EqualFunc(t, 1, 1, equal, msg...) }, false},
		{"EqualFunc/fail", func(t TestingT, msg ...interface{}) { EqualFunc(t, 1, 2, equal, msg...) }, true},
		{"NotEqualFunc/pass", func(t TestingT, msg ...interface{}) { NotEqualFunc(t, 1, 2, equal, msg...) }, false},
		{"NotEqualFunc/fail", func(t TestingT, msg ...interface{}) { NotEqualFunc(t, 1, 1, equal, msg...) }, true},
		{"Less/pass", func(t TestingT, msg ...interface{}) { Less(t, 1, 2, msg...) }, false},
		{"Less/fail", func(t TestingT, msg ...interface{}) { Less(t, 2, 2, msg...) }, true},
		{"LessOrEqual/pass", func(t TestingT, msg ...interface{}) { LessOrEqual(t, 2, 2, msg...) }, false},
		{"LessOrEqual/fail", func(t TestingT, msg ...interface{}) { LessOrEqual(t, 3, 2, msg...) }, true},
		{"Greater/pass", func(t TestingT, msg ...interface{}) { Greater(t, 2, 1, msg...) }, false},
		{"Greater/fail", func(t TestingT, msg ...interface{}) { Greater(t, 2, 2, msg...) }, true},
		{"GreaterOrEqual/pass", func(t TestingT, msg ...interface{}) { GreaterOrEqual(t, 2, 2, msg...) }, false},
		{"GreaterOrEqual/fail", func(t TestingT, msg ...interface{}) { GreaterOrEqual(t, 1, 2, msg...) }, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spy := &testingTSpy{T: t}
			reachedAfter := false
			done := make(chan struct{})
			go func() {
				defer close(done)
				tc.assertion(spy, "case %d", 7)
				reachedAfter = true
			}()
			<-done

			if spy.failCalled != tc.wantFail {
				t.Errorf("failCalled = %v, want %v", spy.failCalled, tc.wantFail)
			}
			if spy.failNowCalled {
				t.Errorf("failNowCalled = true, want false")
			}
			if spy.logCalled != tc.wantFail {
				t.Errorf("logCalled = %v, want %v", spy.logCalled, tc.wantFail)
			}
			if tc.wantFail && !strings.HasPrefix(spy.logMessage, "case 7\n") {
				t.Errorf("logMessage = %q, want prefix %q", spy.logMessage, "case 7\n")
			}
			if !reachedAfter {
				t.Errorf("reachedAfter = false, want true")
			}
		})
	}
}
