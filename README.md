# go-lazy

A Go library for lazy evaluation of values with generic type support.

## Overview

`go-lazy` provides a `Value[T]` type that can hold either an immediate value or a lazy function that will be evaluated when accessed. This is useful for deferring expensive computations until they're actually needed.

## Features

- **Generic type support**: Works with any Go type using generics
- **Immediate values**: Store and retrieve values directly
- **Lazy evaluation**: Defer computation until the value is accessed
- **Type-safe**: Full type safety with Go generics

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