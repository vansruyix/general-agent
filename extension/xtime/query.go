// Package xtime
// @author: fengyi
// @date: 2024/8/12
// @note:
package xtime

// SplitTimeBlock
// Description: 按步长分割时间
// param start 开始时间戳
// param end 结束时间戳
// param size 大小
func SplitTimeBlock(start, end, size int64) [][]int64 {
	step := (end - start) / size
	var times [][]int64
	for s := start; s < end; s += step {
		times = append(times, []int64{s, s + step})
	}
	return times
}
