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

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 || (len(mock.args) == 3 && (mock.args[0] != "" || mock.args[1] != 42 || mock.args[2] != 43)) {
			t.Errorf("Expected args %v, %v and %v, got: %v", "", 42, 43, mock.args)
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

		if mock.format != "%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != "" || mock.args[1] != 42)) {
			t.Errorf("Expected args %v and %v, got: %v", "", 42, mock.args)
		}
	})
}

func TestEqualFunc(t *testing.T) {
	t.Run("passes when custom equality function returns true", func(t *testing.T) {
		mock := &mockT{}
		type Person struct {
			Name string
			Age  int
		}
		equalByName := func(a, b Person) bool {
			return a.Name == b.Name
		}
		EqualFunc(mock, Person{Name: "Alice", Age: 30}, Person{Name: "Alice", Age: 25}, equalByName)
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

	t.Run("fails when custom equality function returns false", func(t *testing.T) {
		mock := &mockT{}
		type Person struct {
			Name string
			Age  int
		}
		equalByName := func(a, b Person) bool {
			return a.Name == b.Name
		}
		expected := Person{Name: "Alice", Age: 30}
		actual := Person{Name: "Bob", Age: 30}
		EqualFunc(mock, expected, actual, equalByName)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 || (len(mock.args) == 3 && (mock.args[0] != "" || mock.args[1] != expected || mock.args[2] != actual)) {
			t.Errorf("Expected args %v, %v and %v, got: %v", "", expected, actual, mock.args)
		}
	})

	t.Run("works with slices using custom equality", func(t *testing.T) {
		mock := &mockT{}
		sliceEqual := func(a, b []int) bool {
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		}
		EqualFunc(mock, []int{1, 2, 3}, []int{1, 2, 3}, sliceEqual)
		if mock.failed {
			t.Error("Expected assertion to pass")
		}
	})
}

func TestNotEqualFunc(t *testing.T) {
	t.Run("passes when custom equality function returns false", func(t *testing.T) {
		mock := &mockT{}
		type Person struct {
			Name string
			Age  int
		}
		equalByName := func(a, b Person) bool {
			return a.Name == b.Name
		}
		NotEqualFunc(mock, Person{Name: "Alice", Age: 30}, Person{Name: "Bob", Age: 30}, equalByName)
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

	t.Run("fails when custom equality function returns true", func(t *testing.T) {
		mock := &mockT{}
		type Person struct {
			Name string
			Age  int
		}
		equalByName := func(a, b Person) bool {
			return a.Name == b.Name
		}
		person := Person{Name: "Alice", Age: 30}
		NotEqualFunc(mock, person, Person{Name: "Alice", Age: 25}, equalByName)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\nvalues to be different\n-----got:----------------------------\nboth are %v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 2 || (len(mock.args) == 2 && (mock.args[0] != "" || mock.args[1] != person)) {
			t.Errorf("Expected args %v and %v, got: %v", "", person, mock.args)
		}
	})

	t.Run("works with slices using custom equality", func(t *testing.T) {
		mock := &mockT{}
		sliceEqual := func(a, b []int) bool {
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		}
		NotEqualFunc(mock, []int{1, 2, 3}, []int{1, 2, 4}, sliceEqual)
		if mock.failed {
			t.Error("Expected assertion to pass")
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

		if mock.format != "%s\nexpected:----------------------------\n%v to be less than %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 43 || mock.args[2] != 42 || mock.args[3] != 43)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 43, 42, 43, mock.args)
		}
	})

	t.Run("fails when actual is equal to max", func(t *testing.T) {
		mock := &mockT{}
		Less(mock, 42, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v to be less than %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 42 || mock.args[2] != 42 || mock.args[3] != 42)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 42, 42, 42, mock.args)
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

		if mock.format != "%s\nexpected:----------------------------\n%v to be less than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 43 || mock.args[2] != 42 || mock.args[3] != 43)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 43, 42, 43, mock.args)
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

		if mock.format != "%s\nexpected:----------------------------\n%v to be greater than %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 42 || mock.args[2] != 43 || mock.args[3] != 42)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 42, 43, 42, mock.args)
		}
	})

	t.Run("fails when actual is equal to min", func(t *testing.T) {
		mock := &mockT{}
		Greater(mock, 42, 42)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v to be greater than %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 42 || mock.args[2] != 42 || mock.args[3] != 42)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 42, 42, 42, mock.args)
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

		if mock.format != "%s\nexpected:----------------------------\n%v to be greater than or equal to %v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 4 || (len(mock.args) == 4 && (mock.args[0] != "" || mock.args[1] != 42 || mock.args[2] != 43 || mock.args[3] != 42)) {
			t.Errorf("Expected args %v, %v, %v and %v, got: %v", "", 42, 43, 42, mock.args)
		}
	})
}

func TestIsNil(t *testing.T) {
	t.Run("passes when value is nil", func(t *testing.T) {
		mock := &mockT{}
		IsNil(mock, nil)
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

	t.Run("passes when pointer is nil", func(t *testing.T) {
		mock := &mockT{}
		var ptr *int
		IsNil(mock, ptr)
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

	t.Run("passes when slice is nil", func(t *testing.T) {
		mock := &mockT{}
		var slice []int
		IsNil(mock, slice)
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

	t.Run("passes when map is nil", func(t *testing.T) {
		mock := &mockT{}
		var m map[string]int
		IsNil(mock, m)
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

	t.Run("passes when channel is nil", func(t *testing.T) {
		mock := &mockT{}
		var ch chan int
		IsNil(mock, ch)
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

	t.Run("passes when interface is nil", func(t *testing.T) {
		mock := &mockT{}
		var iface interface{}
		IsNil(mock, iface)
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

	t.Run("fails when value is not nil", func(t *testing.T) {
		mock := &mockT{}
		value := 42
		IsNil(mock, &value)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when slice is not nil", func(t *testing.T) {
		mock := &mockT{}
		slice := []int{1, 2, 3}
		IsNil(mock, slice)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("passes when function is nil", func(t *testing.T) {
		mock := &mockT{}
		var fn func()
		IsNil(mock, fn)
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

	t.Run("fails when function is not nil", func(t *testing.T) {
		mock := &mockT{}
		fn := func() {}
		IsNil(mock, fn)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})
}

func TestIsNotNil(t *testing.T) {
	t.Run("passes when value is not nil", func(t *testing.T) {
		mock := &mockT{}
		value := 42
		IsNotNil(mock, &value)
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

	t.Run("passes when slice is not nil", func(t *testing.T) {
		mock := &mockT{}
		slice := []int{1, 2, 3}
		IsNotNil(mock, slice)
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

	t.Run("passes when map is not nil", func(t *testing.T) {
		mock := &mockT{}
		m := make(map[string]int)
		IsNotNil(mock, m)
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

	t.Run("passes when channel is not nil", func(t *testing.T) {
		mock := &mockT{}
		ch := make(chan int)
		IsNotNil(mock, ch)
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

	t.Run("fails when value is nil", func(t *testing.T) {
		mock := &mockT{}
		IsNotNil(mock, nil)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when pointer is nil", func(t *testing.T) {
		mock := &mockT{}
		var ptr *int
		IsNotNil(mock, ptr)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when slice is nil", func(t *testing.T) {
		mock := &mockT{}
		var slice []int
		IsNotNil(mock, slice)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when map is nil", func(t *testing.T) {
		mock := &mockT{}
		var m map[string]int
		IsNotNil(mock, m)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when channel is nil", func(t *testing.T) {
		mock := &mockT{}
		var ch chan int
		IsNotNil(mock, ch)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("fails when interface is nil", func(t *testing.T) {
		mock := &mockT{}
		var iface interface{}
		IsNotNil(mock, iface)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})

	t.Run("passes when function is not nil", func(t *testing.T) {
		mock := &mockT{}
		fn := func() {}
		IsNotNil(mock, fn)
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

	t.Run("fails when function is nil", func(t *testing.T) {
		mock := &mockT{}
		var fn func()
		IsNotNil(mock, fn)
		if !mock.failed {
			t.Error("Expected assertion to fail")
		}

		if mock.format != "%s\nexpected:----------------------------\n%v\n-----got:----------------------------\n%v\n-------------------------------------" {
			t.Errorf("Expected assertion message, got: %s", mock.format)
		}

		if len(mock.args) != 3 {
			t.Errorf("Expected 3 args, got: %v", mock.args)
		}
	})
}
