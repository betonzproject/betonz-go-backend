package sliceutils

func Map[T, U any](ts []T, f func(t T) U) []U {
	us := make([]U, len(ts))
	for i, ning := range ts {
		us[i] = f(ning)
	}
	return us
}
