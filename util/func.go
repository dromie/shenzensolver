package util

import (
	A "github.com/IBM/fp-go/array"
)

func Sum() func([]int) int {
	return func(arr []int) int {
		return A.Reduce(func(X int, Y int) int { return X + Y }, 0)(arr)
	}
}

func Count[AS ~[]T1, PRED ~func(T1) bool, T1 any](pred PRED) func(AS) int {
	return func(arr AS) int {
		return Sum()(A.Map(func(a T1) int {
			if pred(a) {
				return 1
			} else {
				return 0
			}
		})(arr))
	}
}
