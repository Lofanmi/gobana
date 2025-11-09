package gotil

func OrDefault[T comparable](args ...T) (res T) {
	var empty T
	for _, item := range args {
		if item != empty {
			return item
		}
	}
	return empty
}

func OrSliceDefault[T comparable](args ...[]T) (res []T) {
	for _, slice := range args {
		if len(slice) > 0 {
			return slice
		}
	}
	return
}
