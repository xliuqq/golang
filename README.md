# GoLang 语言学习

## 并发

- volatile.go: [一写多读必须用atomic](./concurrency/main/volatile.go)
- TestOnce: [Once中为什么要采用原子性操作的校验？](./concurrency/once_test.go)
- TestMutex_DoubleLock: [Mutex 不支持重入导致死锁](./concurrency/mutex_test.go)
- TestRWMutex_DeadLock: [两次读锁间其它协程写锁导致死锁](./concurrency/mutex_test.go)
- TestMutexChannel_DeadLock: [Mutex跟Channel同时使用的死锁](./concurrency/mutex_test.go)
- TestUnbufferedChannel_GoRoutineLeak: [不带缓冲的channel导致go routine泄露](./concurrency/channel_test.go)
- TestContext_Blocking: [context改变导致goroutine卡住](./concurrency/channel_test.go)
- TestLocalSharedVariable: [GoRoutines共享闭包变量期望值错误](./concurrency/goroutine_test.go)
- TestTimer_Zero: [值为0的Timer.C会立即触发select执行 ](./concurrency/channel_test.go)

## Json
- TestJson_Unmarshal: [通过反射自定义Json序列化和反序列化](./json/json_test.go)


## Yaml
- TestYaml_Unmarshal: [通过反射自定义Yaml序列化和反序列化](./yaml/yaml_test.go)

## Context
- TestCOntext_Cancel: [父Context执行Cancel会执行子Context的cancel](./context/context_test.go)
- TestCOntext_Timeout: [Timeout 超时的用法](./context/context_test.go)