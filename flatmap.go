package lazy

func FlatMap[T any, R any](v Value[T], f func(T) Value[R]) Value[R] {
	return NewLazy(func() R {
		return f(v.Get()).Get()
	})
}
