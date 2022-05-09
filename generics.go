package gocbac

// convert slice to the map where each element of a slice is the key and values is bool
func SliceToBoolMap[S ~[]K, K comparable](slice S) map[K]bool {
	res := map[K]bool{}

	for _, item := range slice {
		res[item] = true
	}

	return res
}

// fill map with same values
func MapFill[M ~map[K]V, K comparable, V any](dict M, keys []K, val V) M {
	for _, key := range keys {
		dict[key] = val
	}

	return dict
}
