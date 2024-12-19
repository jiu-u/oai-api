package array

// Map 函数：应用一个转换函数到切片的每个元素
func Map[T any, R any](input []T, fn func(T) R) []R {
	var result []R
	for _, item := range input {
		result = append(result, fn(item))
	}
	return result
}

// Filter 函数：返回满足条件的元素组成的切片
func Filter[T any](input []T, fn func(T) bool) []T {
	var result []T
	for _, item := range input {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Find 函数：返回第一个满足条件的元素，找不到返回零值
func Find[T any](input []T, fn func(T) bool) (T, bool) {
	var zeroValue T // 泛型类型的零值
	for _, item := range input {
		if fn(item) {
			return item, true
		}
	}
	return zeroValue, false
}

// Reduce 函数：累积切片中的元素
func Reduce[T any, R any](input []T, fn func(R, T) R, initial R) R {
	result := initial
	for _, item := range input {
		result = fn(result, item)
	}
	return result
}

// Some 函数：检查切片中是否有任何元素满足条件
func Some[T any](input []T, fn func(T) bool) bool {
	for _, item := range input {
		if fn(item) {
			return true
		}
	}
	return false
}

// Every 函数：检查切片中的所有元素是否都满足条件
func Every[T any](input []T, fn func(T) bool) bool {
	for _, item := range input {
		if !fn(item) {
			return false
		}
	}
	return true
}

// IndexOf 函数：返回指定元素的索引，找不到返回 -1
func IndexOf[T comparable](input []T, value T) int {
	for i, item := range input {
		if item == value {
			return i
		}
	}
	return -1
}

// LastIndexOf 函数：返回最后一次出现的指定元素的索引，找不到返回 -1
func LastIndexOf[T comparable](input []T, value T) int {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == value {
			return i
		}
	}
	return -1
}
