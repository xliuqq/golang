package concurrency

import (
	"testing"
)

// goroutine 中共享的变量（闭包）被外部修改时，其期望值跟实际值会不一致
func TestLocalSharedVariable(t *testing.T) {
	for i := 17; i <= 21; i++ { // write
		go func(index int) {
			if i != index { // read
				t.Errorf("index [%d] is not the same as local variable [%d]", index, i)
			}
		}(i)
	}
}
