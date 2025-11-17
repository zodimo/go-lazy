package lazy

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 42},
		{"string", "hello"},
		{"bool", true},
		{"float64", 3.14},
		{"zero int", 0},
		{"zero string", ""},
		{"zero bool", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				val := New(v)
				if got := val.Get(); got != v {
					t.Errorf("New(%v).Get() = %v, want %v", v, got, v)
				}
			case string:
				val := New(v)
				if got := val.Get(); got != v {
					t.Errorf("New(%v).Get() = %v, want %v", v, got, v)
				}
			case bool:
				val := New(v)
				if got := val.Get(); got != v {
					t.Errorf("New(%v).Get() = %v, want %v", v, got, v)
				}
			case float64:
				val := New(v)
				if got := val.Get(); got != v {
					t.Errorf("New(%v).Get() = %v, want %v", v, got, v)
				}
			}
		})
	}
}

func TestNewLazy(t *testing.T) {
	t.Run("lazy int", func(t *testing.T) {
		called := false
		val := NewLazy(func() int {
			called = true
			return 42
		})

		if called {
			t.Error("Lazy function should not be called during NewLazy")
		}

		result := val.Get()
		if !called {
			t.Error("Lazy function should be called during Get")
		}
		if result != 42 {
			t.Errorf("Get() = %v, want 42", result)
		}
	})

	t.Run("lazy string", func(t *testing.T) {
		val := NewLazy(func() string {
			return "lazy value"
		})

		if got := val.Get(); got != "lazy value" {
			t.Errorf("Get() = %v, want 'lazy value'", got)
		}
	})

	t.Run("lazy function called multiple times", func(t *testing.T) {
		callCount := 0
		val := NewLazy(func() int {
			callCount++
			return callCount * 10
		})

		// Each Get() should call the lazy function
		if got := val.Get(); got != 10 {
			t.Errorf("First Get() = %v, want 10", got)
		}
		if got := val.Get(); got != 20 {
			t.Errorf("Second Get() = %v, want 20", got)
		}
		if callCount != 2 {
			t.Errorf("Lazy function called %d times, want 2", callCount)
		}
	})

	t.Run("lazy with zero value", func(t *testing.T) {
		val := NewLazy(func() int {
			return 0
		})

		if got := val.Get(); got != 0 {
			t.Errorf("Get() = %v, want 0", got)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("immediate value", func(t *testing.T) {
		val := New(100)
		if got := val.Get(); got != 100 {
			t.Errorf("Get() = %v, want 100", got)
		}
		// Get should be idempotent for immediate values
		if got := val.Get(); got != 100 {
			t.Errorf("Second Get() = %v, want 100", got)
		}
	})

	t.Run("lazy value", func(t *testing.T) {
		val := NewLazy(func() string {
			return "computed"
		})
		if got := val.Get(); got != "computed" {
			t.Errorf("Get() = %v, want 'computed'", got)
		}
	})

	t.Run("zero value struct", func(t *testing.T) {
		var val Value[int]
		if got := val.Get(); got != 0 {
			t.Errorf("Get() on zero Value = %v, want 0", got)
		}
	})

	t.Run("zero value string", func(t *testing.T) {
		var val Value[string]
		if got := val.Get(); got != "" {
			t.Errorf("Get() on zero Value = %v, want empty string", got)
		}
	})
}

func TestValueTypes(t *testing.T) {
	t.Run("struct type", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		p := Person{Name: "Alice", Age: 30}
		val := New(p)
		got := val.Get()

		if got.Name != "Alice" || got.Age != 30 {
			t.Errorf("Get() = %+v, want {Name: Alice, Age: 30}", got)
		}
	})

	t.Run("pointer type", func(t *testing.T) {
		x := 42
		val := New(&x)
		got := val.Get()

		if got == nil {
			t.Error("Get() returned nil pointer")
		}
		if *got != 42 {
			t.Errorf("Get() = %v, want 42", *got)
		}
	})

	t.Run("slice type", func(t *testing.T) {
		slice := []int{1, 2, 3}
		val := New(slice)
		got := val.Get()

		if len(got) != 3 {
			t.Errorf("Get() returned slice of length %d, want 3", len(got))
		}
		if got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Errorf("Get() = %v, want [1, 2, 3]", got)
		}
	})

	t.Run("map type", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		val := New(m)
		got := val.Get()

		if len(got) != 2 {
			t.Errorf("Get() returned map of length %d, want 2", len(got))
		}
		if got["a"] != 1 || got["b"] != 2 {
			t.Errorf("Get() = %v, want map[a:1 b:2]", got)
		}
	})
}

func TestLazyEvaluation(t *testing.T) {
	t.Run("lazy function with side effects", func(t *testing.T) {
		sideEffect := 0
		val := NewLazy(func() int {
			sideEffect++
			return sideEffect
		})

		// Side effect should not occur until Get() is called
		if sideEffect != 0 {
			t.Errorf("Side effect occurred before Get(), got %d, want 0", sideEffect)
		}

		result1 := val.Get()
		if sideEffect != 1 {
			t.Errorf("Side effect should be 1 after first Get(), got %d", sideEffect)
		}
		if result1 != 1 {
			t.Errorf("First Get() = %v, want 1", result1)
		}

		result2 := val.Get()
		if sideEffect != 2 {
			t.Errorf("Side effect should be 2 after second Get(), got %d", sideEffect)
		}
		if result2 != 2 {
			t.Errorf("Second Get() = %v, want 2", result2)
		}
	})

	t.Run("lazy function with closure", func(t *testing.T) {
		x := 10
		val := NewLazy(func() int {
			return x * 2
		})

		if got := val.Get(); got != 20 {
			t.Errorf("Get() = %v, want 20", got)
		}

		// Modify closure variable
		x = 20
		if got := val.Get(); got != 40 {
			t.Errorf("Get() after closure change = %v, want 40", got)
		}
	})
}
