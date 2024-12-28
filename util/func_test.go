package util

import (
	"testing"
)

func Test_Sum(t *testing.T) {
	sum := Sum()
	if sum([]int{1, 2, 3}) != 6 {
		t.Errorf("Sum is not 6")
	}
}

func Test_Count(t *testing.T) {
	count := Count[[]int](func(a int) bool { return a%2 == 0 })
	if count([]int{1, 2, 3, 4}) != 2 {
		t.Errorf("Count is not 2")
	}
}
