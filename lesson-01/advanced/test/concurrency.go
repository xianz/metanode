package main

import (
	"fmt"
	"sync"
	"time"
)

/**
并发：线程安全（操作变量时，不会因为多个线程原因造成数据紊乱）
通过加锁来保证线程安全，保证数据操作原子性
sync.Mutex   .Lock()

sync.RMutex  .Unlock()

// 等待执行
sync.WaitGroup
*/

func main() {
	// mutexDemo()
	RWMutexDemo()
}

type safeCounter struct {
	mu    sync.Mutex
	count int
}

// 保证数据原子性
func mutexDemo() {
	fmt.Println("=== mutex ===")
	scounter := &safeCounter{}

	for i := 0; i < 10; i++ {
		go func(id int) {
			scounter.mu.Lock()
			scounter.count++
			fmt.Println("goroutine_A", id, "count:", scounter.count)
			scounter.mu.Unlock()
		}(i)
	}

	for i := 0; i < 10; i++ {
		go func(id int) {
			scounter.mu.Lock()
			scounter.count++
			fmt.Println("goroutine_B", id, " count:", scounter.count)
			scounter.mu.Unlock()
		}(i)
	}
	fmt.Scanln()
}

// 读写分离，写时不能读，读写
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func (sm *SafeMap) write(label string, key string, value int, hold time.Duration) {
	fmt.Printf("[%s] 准备写入 %s\n", label, key)
	sm.mu.Lock()
	fmt.Printf("[%s] 获得写锁，开始写入\n", label)
	sm.data["a"]++
	if hold > 0 {
		time.Sleep(hold)
	}
	fmt.Printf("[%s] 写入完成，释放写锁\n", label)
	sm.mu.Unlock()
}

func (sm *SafeMap) read(label string, key string, hold time.Duration) {
	fmt.Printf("[%s] 等待读\n", label)
	sm.mu.RLock()
	fmt.Printf("[%s] 获得读锁\n", label)
	value, ok := sm.data[key]
	if !ok {
		fmt.Printf("[%s] 还没写[%s]\n", label, key)
	} else {
		fmt.Printf("[%s] %s = %v\n", label, key, value)
	}
	sm.mu.RUnlock()
}

// 线程安全
// 更多解决方式：https://chat.deepseek.com/a/chat/s/f2b66b63-8493-48a6-8fee-cf2af7f38c26
func RWMutexDemo() {
	sm := &SafeMap{data: make(map[string]int)}

	// // 读多写少场景
	go func() {
		sm.write("writer#1", "a", 1, time.Millisecond*500)
	}()
	time.Sleep(time.Second)

	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		sm.read("reader#", "a", 0)
	// 	}()
	// }

	// 读场景
	for i := 0; i < 8; i++ {
		go func(id int) {
			sm.read("only read", "a", 0)
		}(i)
	}

	fmt.Scanln()
}
