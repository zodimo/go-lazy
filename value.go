package lazy

type wrapper[T any] struct {
	value T
}

func (w *wrapper[T]) Get() T {
	return w.value
}

type Value[T any] struct {
	wrapper *wrapper[T]
	lazy    func() T
	isLazy  bool
}

func New[T any](value T) Value[T] {
	return Value[T]{
		wrapper: &wrapper[T]{
			value: value,
		},
		isLazy: false,
	}
}

func NewLazy[T any](lazy func() T) Value[T] {
	return Value[T]{
		lazy:   lazy,
		isLazy: true,
	}
}

func (l Value[T]) Get() T {
	if l.isLazy {
		return l.lazy()
	}
	if l.wrapper == nil {
		zero := new(T)
		return *zero
	}
	return l.wrapper.Get()
}
