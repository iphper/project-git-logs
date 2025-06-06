package utils

// @func Slice_Duplicates 去除切片中的重复元素
func Slice_Duplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	index := 0

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			slice[index] = item
			index++
		}
	}
	return slice[:index]
}
