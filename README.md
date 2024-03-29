# GoLang 语言学习

## 并发

- TestOnce: [Once中为什么要采用原子性操作的校验？](./concurrency/once_test.go)
- TestMutex_DoubleLock: [Mutex 不支持重入导致死锁](./concurrency/mutex_test.go)
- TestRWMutex_DeadLock: [两次读锁间其它协程写锁导致死锁](./concurrency/mutex_test.go)
- TestMutexChannel_DeadLock: [Mutex跟Channel同时使用的死锁](./concurrency/mutex_test.go)
- TestUnbufferedChannel_GoRoutineLeak: [不带缓冲的channel导致go routine泄露](./concurrency/channel_test.go)
- TestContext_Blocking: [context改变导致goroutine卡住](./concurrency/channel_test.go)
- TestLocalSharedVariable: [GoRoutines共享闭包变量期望值错误](./concurrency/goroutine_test.go)
- TestTimer_Zero: [值为0的Timer.C会立即触发select执行 ](./concurrency/channel_test.go)

## Context
- [context 用法示例](./ctx/ctx.go)
  - WithValue 用法
  - WithTimeout 用法（含 Cancel 用法）

## HTTP 
- sse 示例
  - [sse server](./http/server.go)
  - [sse client](./http/client.go)
- websocket 示例
  - [websocket server](./http/ws_server.go)
  - [websocket client](./http/ws_client.go)