package ptrutil

func ValuePtr[T any, V interface {
	Value() T
}](value *V) *T {
	if value == nil {
		return nil
	}

	return new((*value).Value())
}
