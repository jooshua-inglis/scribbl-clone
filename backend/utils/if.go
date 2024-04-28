package utils

func If[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}
