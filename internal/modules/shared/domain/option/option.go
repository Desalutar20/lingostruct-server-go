package option

type Option[T any] struct {
	value T
	valid bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		valid: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsSome() bool {
	return o.valid
}

func (o Option[T]) IsNone() bool {
	return !o.valid
}

func (o Option[T]) Unwrap() T {
	return o.value
}

func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.valid {
		return o.value
	}

	return defaultValue
}
