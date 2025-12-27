package slice_helpers

func Filter[T any](values []T, predicate func(val T) bool) (result []T) {
	for _, val := range values {
		if predicate(val) {
			result = append(result, val)
		}
	}

	return
}
