package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// context 可以作为请求的上下文进行传递，但是需要显示作为函数参数传递
func ctxValueUsage() {
	ctx := context.WithValue(context.TODO(), "a", "b")
	newCtx := context.WithValue(ctx, "c", "d")

	// print b and d
	fmt.Println(newCtx.Value("a"))
	fmt.Println(newCtx.Value("c"))
}

// cancel 可以用来控制同步多个 goroutine，例如取消定时执行的 goroutine
func ctxCancelUsage() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 子 goroutine 需要通知结束
	go func(i int) {
		if i == 1 {
			cancel()
		}
	}(1)

	// 等待任务完成，如果被取消，ctx.Err 会返回 Canceled
	select {
	case <-ctx.Done():
		fmt.Println("main canceled", errors.Is(ctx.Err(), context.Canceled))
	}
}

// 超时的使用场景（等价于 Deadline），例如控制长时间处理的 goroutine 终止
func ctxTimeoutUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 执行任务
	go handle(ctx)

	// 等待任务完成，如果超时，ctx.Err 会返回 deadline
	select {
	case <-ctx.Done():
		fmt.Println("main deadline exceeded", errors.Is(ctx.Err(), context.DeadlineExceeded))
	}
}

func handle(ctx context.Context) {
	// 一直执行任务
	for {
		select {
		default:
			fmt.Println("work")
			time.Sleep(300 * time.Millisecond)
		case <-ctx.Done():
			fmt.Println("work done")
			return
		}
	}
}

func main() {
	ctxValueUsage()
	//ctxTimeoutUsage()
	//ctxCancelUsage()
}
