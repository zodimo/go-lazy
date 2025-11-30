# go-lazy

A Go library for lazy evaluation of values with generic type support.

## Overview

`go-lazy` provides a `Value[T]` type that can hold either an immediate value or a lazy function that will be evaluated when accessed. This is useful for deferring expensive computations until they're actually needed.

## Features

- **Generic type support**: Works with any Go type using generics
- **Immediate values**: Store and retrieve values directly
- **Lazy evaluation**: Defer computation until the value is accessed
- **Type-safe**: Full type safety with Go generics
- **Map and FlatMap**: Transform and chain lazy values with functional operations

## Installation

```bash
go get github.com/zodimo/go-lazy
```

## Usage

### Immediate Values

Create a `Value` with an immediate value using `New`:

```go
package main

import (
    "fmt"
    "github.com/zodimo/go-lazy"
)

func main() {
    // Create an immediate value
    v := lazy.New(42)
    
    // Get the value
    result := v.Get()
    fmt.Println(result) // Output: 42
}
```

### Lazy Values

Create a `Value` with a lazy function using `NewLazy`. The function will be called each time `Get()` is invoked:

```go
package main

import (
    "fmt"
    "github.com/zodimo/go-lazy"
)

func main() {
    // Create a lazy value
    v := lazy.NewLazy(func() int {
        fmt.Println("Computing expensive value...")
        return 100 * 100
    })
    
    // The function is called when Get() is invoked
    result := v.Get() // Output: Computing expensive value...
    fmt.Println(result) // Output: 10000
    
    // Calling Get() again will call the function again
    result2 := v.Get() // Output: Computing expensive value...
    fmt.Println(result2) // Output: 10000
}
```

### Example: Deferring Expensive Operations

```go
package main

import (
    "fmt"
    "github.com/zodimo/go-lazy"
)

func expensiveComputation() string {
    // Simulate expensive work
    return "Computed result"
}

func main() {
    // Create a lazy value - computation won't happen yet
    result := lazy.NewLazy(expensiveComputation)
    
    // Use the value only if needed
    if someCondition {
        value := result.Get() // Computation happens here
        fmt.Println(value)
    }
    // If someCondition is false, expensiveComputation is never called
}
```

### Map

Transform a `Value` by applying a function to it. The transformation is lazy and only happens when the result is accessed:

```go
package main

import (
    "fmt"
    "github.com/zodimo/go-lazy"
)

func main() {
    // Create a value
    val := lazy.New(5)
    
    // Map it to double the value
    doubled := lazy.Map(val, func(x int) int {
        return x * 2
    })
    
    // The transformation happens lazily when Get() is called
    result := doubled.Get()
    fmt.Println(result) // Output: 10
    
    // Map can also change types
    length := lazy.Map(lazy.New("hello"), func(s string) int {
        return len(s)
    })
    fmt.Println(length.Get()) // Output: 5
}
```

### FlatMap

Transform a `Value` by applying a function that returns another `Value`, then flattening the result. This is useful for chaining operations:

```go
package main

import (
    "fmt"
    "github.com/zodimo/go-lazy"
)

func main() {
    // Create a value
    val := lazy.New(3)
    
    // FlatMap with a function that returns a Value
    result := lazy.FlatMap(val, func(x int) lazy.Value[int] {
        return lazy.New(x * 2)
    })
    
    fmt.Println(result.Get()) // Output: 6
    
    // FlatMap can chain with lazy values
    chained := lazy.FlatMap(val, func(x int) lazy.Value[int] {
        return lazy.NewLazy(func() int {
            return x * 4
        })
    })
    
    fmt.Println(chained.Get()) // Output: 12
    
    // You can chain Map and FlatMap together
    mapped := lazy.Map(val, func(x int) int {
        return x * 2
    })
    flatMapped := lazy.FlatMap(mapped, func(x int) lazy.Value[int] {
        return lazy.New(x + 1)
    })
    fmt.Println(flatMapped.Get()) // Output: 7
}
```

## API Reference

### Types

#### `Value[T any]`

A generic type that holds either an immediate value or a lazy function.

### Functions

#### `New[T any](value T) Value[T]`

Creates a new `Value` with an immediate value.

**Parameters:**
- `value`: The value to store

**Returns:**
- `Value[T]`: A new Value containing the immediate value

#### `NewLazy[T any](lazy func() T) Value[T]`

Creates a new `Value` with a lazy function that will be evaluated when `Get()` is called.

**Parameters:**
- `lazy`: A function that returns a value of type `T`

**Returns:**
- `Value[T]`: A new Value that will evaluate the lazy function on demand

**Note:** The lazy function is called every time `Get()` is invoked. If you need memoization (evaluation once and caching), you'll need to implement that yourself.

#### `Map[T any, R any](v Value[T], f func(T) R) Value[R]`

Transforms a `Value[T]` into a `Value[R]` by applying a function to the value. The transformation is lazy and only occurs when the result is accessed.

**Parameters:**
- `v`: The source `Value[T]` to transform
- `f`: A function that takes a value of type `T` and returns a value of type `R`

**Returns:**
- `Value[R]`: A new lazy Value that will apply the transformation when accessed

**Note:** The transformation function is called lazily when `Get()` is invoked on the returned Value. If the source Value is lazy, both the source evaluation and the transformation happen on the same `Get()` call.

#### `FlatMap[T any, R any](v Value[T], f func(T) Value[R]) Value[R]`

Transforms a `Value[T]` into a `Value[R]` by applying a function that returns a `Value[R]`, then flattening the result. This is useful for chaining operations where each step returns a lazy Value.

**Parameters:**
- `v`: The source `Value[T]` to transform
- `f`: A function that takes a value of type `T` and returns a `Value[R]`

**Returns:**
- `Value[R]`: A new lazy Value that will apply the transformation and flatten the result when accessed

**Note:** The transformation function is called lazily when `Get()` is invoked on the returned Value. The function can return either an immediate Value (using `New`) or a lazy Value (using `NewLazy`), and both will be handled correctly.

### Methods

#### `(l Value[T]) Get() T`

Retrieves the value. For immediate values, returns the stored value. For lazy values, calls the lazy function and returns its result.

**Returns:**
- `T`: The value (either immediate or computed from the lazy function)

## Notes

- Lazy values are **not memoized** by default. Each call to `Get()` on a lazy value will invoke the lazy function again.
- If you need memoization (evaluate once and cache), consider wrapping the lazy function to cache the result.
- Zero values are returned if a Value is in an invalid state (nil wrapper for immediate values).

## License

MIT License

Copyright (c) 2024

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.