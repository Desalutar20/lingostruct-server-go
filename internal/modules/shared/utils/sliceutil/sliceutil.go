package sliceutil

func Every[T any](slice []T, predicate func(value T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}

	return true
}

func Filter[T any](slice []T, predicate func(value T) bool) []T {
	filtered := make([]T, 0, len(slice))

	for _, v := range slice {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
