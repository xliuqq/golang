package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext_Cancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	parent := func() {
		time.Sleep(3 * time.Second)
		cancel()
		fmt.Println(time.Now().String(), "--", "parent cancel", ctx.Err())
	}

	go func(ctx context.Context) {
		ctx, cancel := context.WithCancel(ctx)
		defer func() {
			fmt.Println(time.Now().String(), "--", "child cancel")
			cancel()
		}()
		fmt.Println(time.Now().String(), "--", "run child")
		select {
		// 等待 ctx 结束，当父 ctx 的 cancel 函数调用时，会执行到这里
		case <-ctx.Done():
			fmt.Println(time.Now().String(), "--", "func finished", ctx.Err())
		}
	}(ctx)

	parent()
}

func TestContext_Timeout(t *testing.T) {
	// timeout 内部会有 time.
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	// cancel 一定要执行
	defer cancel()

	// 子 goroutine 会返回
	go func(ctx context.Context) {
		time.Sleep(500 * time.Millisecond) // 第 1 段
		fmt.Println(time.Now().String(), "--first execute over")
		select {
		case <-ctx.Done():
			fmt.Println(time.Now().String(), "--ctx timeout, return")
			return
		default:
			// go next
		}
		time.Sleep(500 * time.Millisecond) // 第 2 段
		fmt.Println(time.Now().String(), "--second execute over")
		// 执行结束后调用 cancel
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		fmt.Println(time.Now().String(), "--timeout")
	}

	time.Sleep(2 * time.Second)
}
