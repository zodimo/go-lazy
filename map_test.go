package lazy

import (
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("map int to int", func(t *testing.T) {
		val := New(5)
		mapped := Map(val, func(x int) int {
			return x * 2
		})

		if got := mapped.Get(); got != 10 {
			t.Errorf("Map(5, x*2).Get() = %v, want 10", got)
		}
	})

	t.Run("map int to string", func(t *testing.T) {
		val := New(42)
		mapped := Map(val, func(x int) string {
			return "value"
		})

		if got := mapped.Get(); got != "value" {
			t.Errorf("Map(42, func).Get() = %v, want 'value'", got)
		}
	})

	t.Run("map string to int", func(t *testing.T) {
		val := New("hello")
		mapped := Map(val, func(s string) int {
			return len(s)
		})

		if got := mapped.Get(); got != 5 {
			t.Errorf("Map('hello', len).Get() = %v, want 5", got)
		}
	})

	t.Run("map with lazy value", func(t *testing.T) {
		called := false
		val := NewLazy(func() int {
			called = true
			return 7
		})

		mapped := Map(val, func(x int) int {
			return x * 3
		})

		// Lazy evaluation: original value should not be called yet
		if called {
			t.Error("Original lazy function should not be called during Map")
		}

		// Get should trigger both lazy evaluations
		result := mapped.Get()
		if !called {
			t.Error("Original lazy function should be called during Get")
		}
		if result != 21 {
			t.Errorf("Map(lazy(7), x*3).Get() = %v, want 21", result)
		}
	})

	t.Run("map function is lazy", func(t *testing.T) {
		mapCalled := false
		val := New(10)
		mapped := Map(val, func(x int) int {
			mapCalled = true
			return x + 1
		})

		// Map function should not be called during Map creation
		if mapCalled {
			t.Error("Map function should not be called during Map creation")
		}

		// Map function should be called during Get
		result := mapped.Get()
		if !mapCalled {
			t.Error("Map function should be called during Get")
		}
		if result != 11 {
			t.Errorf("Map(10, x+1).Get() = %v, want 11", result)
		}
	})

	t.Run("map zero value", func(t *testing.T) {
		val := New(0)
		mapped := Map(val, func(x int) int {
			return x + 1
		})

		if got := mapped.Get(); got != 1 {
			t.Errorf("Map(0, x+1).Get() = %v, want 1", got)
		}
	})

	t.Run("map multiple times", func(t *testing.T) {
		val := New(2)
		mapped1 := Map(val, func(x int) int {
			return x * 2
		})
		mapped2 := Map(mapped1, func(x int) int {
			return x + 1
		})

		if got := mapped2.Get(); got != 5 {
			t.Errorf("Map(Map(2, x*2), x+1).Get() = %v, want 5", got)
		}
	})

	t.Run("map with struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		val := New(Person{Name: "Alice", Age: 30})
		mapped := Map(val, func(p Person) string {
			return p.Name
		})

		if got := mapped.Get(); got != "Alice" {
			t.Errorf("Map(Person, Name).Get() = %v, want 'Alice'", got)
		}
	})

	t.Run("map with slice", func(t *testing.T) {
		val := New([]int{1, 2, 3})
		mapped := Map(val, func(s []int) int {
			return len(s)
		})

		if got := mapped.Get(); got != 3 {
			t.Errorf("Map([]int{1,2,3}, len).Get() = %v, want 3", got)
		}
	})
}
