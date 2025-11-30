package lazy

func Map[T any, R any](v Value[T], f func(T) R) Value[R] {
	return NewLazy(func() R {
		return f(v.Get())
	})
}
