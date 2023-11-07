package concurrency

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 不带缓冲的channel导致go routine泄露（select多个选择导致）
func TestUnbufferedChannel_GoRoutineLeak(t *testing.T) {
	// ch := make(chan int, 1) 会解决这个问题
	ch := make(chan int)
	go func() {
		result := 1
		time.Sleep(1 * time.Second)
		// block here, goroutine will never quit.
		ch <- result
		fmt.Println("!! never print this line")
	}()

	// select 的多个选择，会导致 <- ch 不执行，没有缓冲导致上面的 goroutine 卡住
	select {
	case r := <-ch:
		fmt.Println(r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Time out")
	}

	for {

	}
}

func TestContext_Blocking(t *testing.T) {
	ch := make(chan int)

	ctx := context.Background()
	timeout := 4 * time.Second

	hctx, hcancel := context.WithCancel(ctx)
	fmt.Println("outer hctx", &hctx, &hcancel)

	// hctx 被重新赋值，导致 goroutine 无法退出
	go func() {
		fmt.Println("inner hctx", &hctx, &hcancel)
		select {
		case <-hctx.Done():
			fmt.Println("!!! never print")
			ch <- 1
		}
	}()

	// 将以下4行移到 go routine 之前，则程序正常
	time.Sleep(1 * time.Second)
	if timeout > 0 {
		hctx, hcancel = context.WithTimeout(ctx, timeout)
	}
	fmt.Println("outer new hctx", &hctx, &hcancel)

	hcancel()

	// block
	<-ch
}

func TestChannel_MultiClose(t *testing.T) {
	ch := make(chan int)

	close(ch)
	close(ch)
}

func TestTimer_Zero(t *testing.T) {
	dur := 0

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		// use this, when dur equals 0, will return immediately
		timer := time.NewTimer(0)
		if dur > 0 {
			timer = time.NewTimer(time.Duration(dur))
		}
		fmt.Println("1: select here")
		select {
		case <-timer.C:
			fmt.Println("1: timer here")
		case <-ctx.Done():
			fmt.Println("1: ctx here")
		}
	}()

	go func() {
		defer wg.Done()

		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		// use this, when dur equals 0, will return when ctx.Done
		var timeout <-chan time.Time = nil
		if dur > 0 {
			timeout = time.NewTimer(time.Duration(dur)).C
		}
		fmt.Println("2: select here")
		select {
		case <-timeout:
			fmt.Println("2: timeout here")
		case <-ctx.Done():
			fmt.Println("2: ctx here")
		}
	}()

	wg.Wait()
}
