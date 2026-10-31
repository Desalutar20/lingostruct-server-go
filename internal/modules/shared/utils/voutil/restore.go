package voutil

import "fmt"

func Restore[T any, VO any](
	value *T,
	constructor func(T) (VO, error),
) *VO {
	if value == nil {
		return nil
	}

	vo, err := constructor(*value)
	if err != nil {
		panic(fmt.Errorf("failed to restore value object: %w", err))
	}

	return &vo
}
