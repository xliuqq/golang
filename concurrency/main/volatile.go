package main

import (
	"sync/atomic"
	"time"
)

// go run -race volatile.go

/*
一读多写的并发场景，必须要用 atomic 操作（类似 Java 的 volatile）
*/

type Count struct {
	num int32
}

func main() {
	ch := make(chan int32, 3)

	count := &Count{
		num: 0,
	}
	go func(count *Count) {
		for {
			select {
			case number := <-ch:
				// -race 分析会有一致性问题（Data Race）
				//count.num = number
				atomic.StoreInt32(&count.num, number)
				println("receive number is", number)
			}
		}
	}(count)

	go func() {
		var i int32 = 1
		for {
			ch <- i
			i++
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		// -race 分析会有一致性问题
		//println(count.num)
		println(atomic.LoadInt32(&(count.num)))
		time.Sleep(1 * time.Second)
	}
}
