// Package logic
// @author: liulun
// @date: 2024/6/25
// @note: if逻辑
package logic

type IfElse[T any] struct {
	b        bool
	yesValue T
	noValue  T
}

func IfThen[T any](b bool, v1 T, v2 T) T {
	if b {
		return v1
	} else {
		return v2
	}
}
func If[T any](b bool) *IfElse[T] {
	f := IfElse[T]{}
	f.b = b
	return &f
}
func (f *IfElse[T]) Then(yesValue T) *IfElse[T] {
	f.yesValue = yesValue
	return f
}
func (f *IfElse[T]) Else(noValue T) *IfElse[T] {
	f.noValue = noValue
	return f
}
func (f *IfElse[T]) Get() T {
	if f.b {
		return f.yesValue
	} else {
		return f.noValue
	}
}
