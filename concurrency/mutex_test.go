package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Mutex 不支持重入，单个 goroutine 上锁多次导致死锁
func TestMutex_DoubleLock(t *testing.T) {
	mutex := sync.Mutex{}

	mutex.Lock()
	fmt.Println("lock 1")
	mutex.Lock()
	fmt.Println("lock 2")

	mutex.Unlock()
	fmt.Println("unlock 1")

	mutex.Unlock()
	fmt.Println("unlock 2")
}

// RWMutex 支持重入（多个读不会阻塞），单个 goroutine 上锁多次不会死锁
func TestRWMutex_DoubleLock(t *testing.T) {
	mutex := sync.RWMutex{}

	mutex.RLock()
	fmt.Println("read lock 1")
	mutex.RLock()
	fmt.Println("read lock 2")

	mutex.RUnlock()
	fmt.Println("read unlock 1")

	mutex.RUnlock()
	fmt.Println("read unlock 2")
}

// RWMutex goroutine-A 两个RLock 间执行 goroutine-B 的 Lock（写锁），导致死锁
func TestRWMutex_DeadLock(t *testing.T) {
	mutex := sync.RWMutex{}

	write := make(chan int)

	go func() {
		fmt.Println("read lock 1")
		mutex.RLock()
		write <- 1
		// wait for write lock
		time.Sleep(500 * time.Millisecond)

		fmt.Println("read unlock 2")
		mutex.RLock()
		fmt.Println("read locking 2")
		mutex.RUnlock()
		fmt.Println("read unlock 2")
	}()

	<-write
	fmt.Println("write lock")
	mutex.Lock()
	fmt.Println("write locking")
	mutex.Unlock()
	fmt.Println("write unlock")

}

// Mutex 跟 Channel 一起使用时，保护不当会导致死锁
func TestMutexChannel_DeadLock(t *testing.T) {
	mutex := sync.Mutex{}
	ch := make(chan int)

	go func() {
		mutex.Lock()
		ch <- 1 // block

		// 可以用 select 的 default 使得 ch 非阻塞，或者用 buffered channel

		//select {
		//case ch <- 1:
		//	fmt.Println("send to ch")
		//default:
		//	fmt.Println("not send, no block")
		//}
		mutex.Unlock()
		fmt.Println("goroutine finish")
	}()

	// wait for goroutine lock
	time.Sleep(100 * time.Millisecond)

	mutex.Lock()
	mutex.Unlock()
	<-ch
}
