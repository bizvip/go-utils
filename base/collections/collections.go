package collections

import (
	"cmp"
	"slices"
)

// Filter 使用泛型过滤切片，返回满足条件的元素
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map 使用泛型将切片中的每个元素转换为另一种类型
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Reduce 使用泛型对切片进行归约操作
func Reduce[T, U any](slice []T, initialValue U, reducer func(U, T) U) U {
	result := initialValue
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Find 查找第一个满足条件的元素
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Contains 检查切片是否包含指定元素（适用于可比较类型）
func Contains[T comparable](slice []T, target T) bool {
	return slices.Contains(slice, target)
}

// Unique 返回包含唯一元素的切片（适用于可比较类型）
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// SortBy 使用自定义键函数对切片进行排序
func SortBy[T any, K cmp.Ordered](slice []T, keyFunc func(T) K) {
	slices.SortFunc(slice, func(a, b T) int {
		return cmp.Compare(keyFunc(a), keyFunc(b))
	})
}

// GroupBy 按指定键对切片元素进行分组
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// Chunk 将切片分割为指定大小的块
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Reverse 反转切片
func Reverse[T any](slice []T) {
	slices.Reverse(slice)
}
