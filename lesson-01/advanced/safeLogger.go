package main

import (
	"fmt"
	"sync"
)

/*
**并发练习：**
- 实现一个并发安全的日志系统
- 使用goroutine异步写入日志
- 使用WaitGroup确保所有日志都被写入
*/

func main() {
	sl := NewSafeLogger(10)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			wg.Done()
			sl.Log(fmt.Sprintf("goroutine %d", i))
		}(i)
	}
	wg.Wait()
	for _, row := range sl.getLogs() {
		fmt.Println(row)
	}
}

type SafeLogger struct {
	mu      sync.RWMutex
	logs    []string
	maxSize int
}

func NewSafeLogger(maxSize int) *SafeLogger {
	return &SafeLogger{logs: make([]string, 0), maxSize: maxSize}
}

func (sl *SafeLogger) Log(msg string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logs = append(sl.logs, msg)
	if len(sl.logs) > sl.maxSize {
		sl.logs = sl.logs[1:]
	}
}

func (sl *SafeLogger) getLogs() []string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	result := make([]string, len(sl.logs))
	copy(result, sl.logs) // 在调用getLogs时就固定当时的数据，否则sl.logs会持续变化
	return result
}
