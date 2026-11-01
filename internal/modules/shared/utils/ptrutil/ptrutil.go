package ptrutil

func To[T any](value T) *T {
	return &value
}

func ValuePtr[T any, V interface {
	Value() T
}](value *V) *T {
	if value == nil {
		return nil
	}

	return To((*value).Value())
}
