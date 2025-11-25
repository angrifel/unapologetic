package assert

import (
	"testing"
)

// mockT is a mock implementation of testing.T for testing assertions
type mockT struct {
	failed bool
	format string
	args   []any
}

func (m *mockT) Helper() {}

func (m *mockT) Errorf(format string, args ...any) {
	m.failed = true
	m.format = format
	m.args = args
}

func (m *mockT) Fatalf(format string, args ...any) {
	m.failed = true
	m.format = format
	m.args = args
}

func TestEqual(t *testing.T) {
	t.Run("passes when values are equal", func(t *testing.T) {
		mock := &mockT{}
		Equal(mock, 42, 42)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when values are not equal", func(t *testing.T) {
		mock := &mockT{}
		Equal(mock, 42, 43)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v, got %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 42 || mock.args[1] != 43)) {
			t.Errorf("Expected args %v and %v, got: %v", 42, 43, mock.args)
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("passes when values are different", func(t *testing.T) {
		mock := &mockT{}
		NotEqual(mock, 42, 43)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when values are equal", func(t *testing.T) {
		mock := &mockT{}
		NotEqual(mock, 42, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected values to be different, but both are %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 1 || (len(mock.args) == 1 && mock.args[0] != 42) {
			t.Errorf("Expected args %v, got: %v", 42, mock.args)
		}
	})
}

func TestLess(t *testing.T) {
	t.Run("passes when actual is less than max", func(t *testing.T) {
		mock := &mockT{}
		Less(mock, 42, 43)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when actual is greater than max", func(t *testing.T) {
		mock := &mockT{}
		Less(mock, 43, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be less than %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 43 || mock.args[1] != 42)) {
			t.Errorf("Expected args %v and %v, got: %v", 43, 42, mock.args)
		}
	})

	t.Run("fails when actual is equal to max", func(t *testing.T) {
		mock := &mockT{}
		Less(mock, 42, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be less than %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 42 || mock.args[1] != 42)) {
			t.Errorf("Expected args %v and %v, got: %v", 42, 42, mock.args)
		}
	})
}

func TestLessOrEqual(t *testing.T) {
	t.Run("passes when actual is less than max", func(t *testing.T) {
		mock := &mockT{}
		LessOrEqual(mock, 42, 43)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("passes when actual is equal to max", func(t *testing.T) {
		mock := &mockT{}
		LessOrEqual(mock, 42, 42)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when actual is greater than max", func(t *testing.T) {
		mock := &mockT{}
		LessOrEqual(mock, 43, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be less than or equal to %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 43 || mock.args[1] != 42)) {
			t.Errorf("Expected args %v and %v, got: %v", 43, 42, mock.args)
		}
	})
}

func TestGreater(t *testing.T) {
	t.Run("passes when actual is greater than min", func(t *testing.T) {
		mock := &mockT{}
		Greater(mock, 43, 42)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when actual is less than min", func(t *testing.T) {
		mock := &mockT{}
		Greater(mock, 42, 43)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be greater than %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 42 || mock.args[1] != 43)) {
			t.Errorf("Expected args %v and %v, got: %v", 42, 43, mock.args)
		}
	})

	t.Run("fails when actual is equal to min", func(t *testing.T) {
		mock := &mockT{}
		Greater(mock, 42, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be greater than %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 42 || mock.args[1] != 42)) {
			t.Errorf("Expected args %v and %v, got: %v", 42, 42, mock.args)
		}
	})
}

func TestGreaterOrEqual(t *testing.T) {
	t.Run("passes when actual is greater than min", func(t *testing.T) {
		mock := &mockT{}
		GreaterOrEqual(mock, 43, 42)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("passes when actual is equal to min", func(t *testing.T) {
		mock := &mockT{}
		GreaterOrEqual(mock, 42, 42)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}

		if mock.format != "" {
			t.Errorf("Expected no assertion message, got: %s", mock.format)
		}

		if len(mock.args) > 0 {
			t.Errorf("Expected no args, got: %v", mock.args)
		}
	})

	t.Run("fails when actual is less than min", func(t *testing.T) {
		mock := &mockT{}
		GreaterOrEqual(mock, 42, 43)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "Expected %v to be greater than or equal to %v" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != 42 || mock.args[1] != 43)) {
			t.Errorf("Expected args %v and %v, got: %v", 42, 43, mock.args)
		}
	})
}
