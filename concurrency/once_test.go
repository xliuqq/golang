package concurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.

	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		//defer func() { o.done = 1 }()
		// 如果不采用原子性操作，可能会被重排序，导致先执行 o.done == 1，再执行 f()，atomic 变量的语义等价于 Java 中的 volatile
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func TestOnce(t *testing.T) {
	once := Once{}

	for i := 0; i < 10; i++ {
		go func(idx int) {
			once.Do(func() {
				time.Sleep(400 * time.Millisecond)
				fmt.Printf("**********init********\n")
			})
			fmt.Printf("======init over %d========\n", idx)
		}(i)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}
