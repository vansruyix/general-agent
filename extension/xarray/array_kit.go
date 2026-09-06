// Package xarray
// @author: liulun
// @date: 2024/5/28
// @note: 集合工具类
package xarray

import "strconv"

func StrIndexOf(array []string, str string) int {
	for i, item := range array {
		if item == str {
			return i
		}
	}
	return -1
}
func IndexOf[T comparable](array []T, a T) int {
	for i, item := range array {
		if item == a {
			return i
		}
	}
	return -1
}
func StrsToInts(array []string) ([]int, error) {
	newArray := make([]int, len(array))
	for i, item := range array {
		r, err := strconv.Atoi(item)
		if err != nil {
			return []int{}, err
		}
		newArray[i] = r
	}
	return newArray, nil
}

// Deduplicate 去重
func Deduplicate[T comparable](array []T) []T {
	if len(array) <= 1 {
		return array
	}

	// 使用map来记录已经出现的元素，提高查找效率
	seen := make(map[T]bool)
	newArray := make([]T, 0)
	for _, item := range array {
		if !seen[item] {
			seen[item] = true
			newArray = append(newArray, item)
		}
	}
	return newArray
}

// DeduplicateAndExclude 去重并排除指定值
func DeduplicateAndExclude[T comparable](array []T, excludes ...T) []T {
	if len(array) == 0 {
		return array
	}

	// 将要排除的值放入map中，提高查找效率
	excludeMap := make(map[T]bool)
	for _, item := range excludes {
		excludeMap[item] = true
	}

	// 使用map来记录已经出现的元素，提高查找效率
	seen := make(map[T]bool)
	newArray := make([]T, 0)
	for _, item := range array {
		// 如果该元素不在排除列表中，且未重复出现，则添加到结果中
		if !excludeMap[item] && !seen[item] {
			seen[item] = true
			newArray = append(newArray, item)
		}
	}
	return newArray
}
