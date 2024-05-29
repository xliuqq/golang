package basic

import (
	"testing"
)

func adder() func(int) int {
	sum := 0
	return func(x int) int { sum += x; return sum }
}

// 测试函数闭包
func TestClosure(t *testing.T) {
	add := adder()
	add(2)
	sum := add(3)
	if sum != 5 {
		t.Errorf("expected 5, got %d", sum)
	}
	newAdd := adder()
	sum = newAdd(2)
	if sum != 2 {
		t.Errorf("expected 2, got %d", sum)
	}
}
