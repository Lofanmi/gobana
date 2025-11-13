package gotil

func IfThen[T any](condition bool, value T) (res T) {
	if condition {
		res = value
	}
	return
}
