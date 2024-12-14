package utils

/*
 * ConvertSliceType converts a slice of type A to another type B
 */
func ConvertSliceType[A, B any](s []A, f func(A) B) []B {
	result := make([]B, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}
	return result
}
