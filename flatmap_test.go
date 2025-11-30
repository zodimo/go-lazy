package lazy

import (
	"testing"
)

func TestFlatMap(t *testing.T) {
	t.Run("flatmap int to int", func(t *testing.T) {
		val := New(5)
		flatMapped := FlatMap(val, func(x int) Value[int] {
			return New(x * 2)
		})

		if got := flatMapped.Get(); got != 10 {
			t.Errorf("FlatMap(5, x*2).Get() = %v, want 10", got)
		}
	})

	t.Run("flatmap int to string", func(t *testing.T) {
		val := New(42)
		flatMapped := FlatMap(val, func(x int) Value[string] {
			return New("value")
		})

		if got := flatMapped.Get(); got != "value" {
			t.Errorf("FlatMap(42, func).Get() = %v, want 'value'", got)
		}
	})

	t.Run("flatmap with lazy value in function", func(t *testing.T) {
		val := New(3)
		flatMapped := FlatMap(val, func(x int) Value[int] {
			return NewLazy(func() int {
				return x * 4
			})
		})

		// The lazy function inside FlatMap should not be called yet
		// But we can't easily test this without modifying the code
		// So we just test the result
		if got := flatMapped.Get(); got != 12 {
			t.Errorf("FlatMap(3, lazy(x*4)).Get() = %v, want 12", got)
		}
	})

	t.Run("flatmap with lazy source value", func(t *testing.T) {
		called := false
		val := NewLazy(func() int {
			called = true
			return 7
		})

		flatMapped := FlatMap(val, func(x int) Value[int] {
			return New(x * 3)
		})

		// Lazy evaluation: original value should not be called yet
		if called {
			t.Error("Original lazy function should not be called during FlatMap")
		}

		// Get should trigger both lazy evaluations
		result := flatMapped.Get()
		if !called {
			t.Error("Original lazy function should be called during Get")
		}
		if result != 21 {
			t.Errorf("FlatMap(lazy(7), x*3).Get() = %v, want 21", result)
		}
	})

	t.Run("flatmap function is lazy", func(t *testing.T) {
		flatMapCalled := false
		val := New(10)
		flatMapped := FlatMap(val, func(x int) Value[int] {
			flatMapCalled = true
			return New(x + 1)
		})

		// FlatMap function should not be called during FlatMap creation
		if flatMapCalled {
			t.Error("FlatMap function should not be called during FlatMap creation")
		}

		// FlatMap function should be called during Get
		result := flatMapped.Get()
		if !flatMapCalled {
			t.Error("FlatMap function should be called during Get")
		}
		if result != 11 {
			t.Errorf("FlatMap(10, x+1).Get() = %v, want 11", result)
		}
	})

	t.Run("flatmap zero value", func(t *testing.T) {
		val := New(0)
		flatMapped := FlatMap(val, func(x int) Value[int] {
			return New(x + 1)
		})

		if got := flatMapped.Get(); got != 1 {
			t.Errorf("FlatMap(0, x+1).Get() = %v, want 1", got)
		}
	})

	t.Run("flatmap multiple times", func(t *testing.T) {
		val := New(2)
		flatMapped1 := FlatMap(val, func(x int) Value[int] {
			return New(x * 2)
		})
		flatMapped2 := FlatMap(flatMapped1, func(x int) Value[int] {
			return New(x + 1)
		})

		if got := flatMapped2.Get(); got != 5 {
			t.Errorf("FlatMap(FlatMap(2, x*2), x+1).Get() = %v, want 5", got)
		}
	})

	t.Run("flatmap chaining map and flatmap", func(t *testing.T) {
		val := New(3)
		mapped := Map(val, func(x int) int {
			return x * 2
		})
		flatMapped := FlatMap(mapped, func(x int) Value[int] {
			return New(x + 1)
		})

		if got := flatMapped.Get(); got != 7 {
			t.Errorf("FlatMap(Map(3, x*2), x+1).Get() = %v, want 7", got)
		}
	})

	t.Run("flatmap with struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		val := New(Person{Name: "Bob", Age: 25})
		flatMapped := FlatMap(val, func(p Person) Value[string] {
			return New(p.Name)
		})

		if got := flatMapped.Get(); got != "Bob" {
			t.Errorf("FlatMap(Person, Name).Get() = %v, want 'Bob'", got)
		}
	})

	t.Run("flatmap returning lazy value", func(t *testing.T) {
		callCount := 0
		val := New(5)
		flatMapped := FlatMap(val, func(x int) Value[int] {
			return NewLazy(func() int {
				callCount++
				return x * 2
			})
		})

		// The lazy function should not be called yet
		if callCount != 0 {
			t.Errorf("Lazy function called %d times, want 0", callCount)
		}

		// First Get should call the lazy function
		if got := flatMapped.Get(); got != 10 {
			t.Errorf("First FlatMap.Get() = %v, want 10", got)
		}
		if callCount != 1 {
			t.Errorf("Lazy function called %d times, want 1", callCount)
		}

		// Second Get should call the lazy function again
		if got := flatMapped.Get(); got != 10 {
			t.Errorf("Second FlatMap.Get() = %v, want 10", got)
		}
		if callCount != 2 {
			t.Errorf("Lazy function called %d times, want 2", callCount)
		}
	})
}
