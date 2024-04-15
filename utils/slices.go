package utils

func SlicesMap[T any, U any](s []T, f func(T) U) []U {
	res := make([]U, 0, len(s))
	for _, v := range s {
		res = append(res, f(v))
	}
	return res
}
