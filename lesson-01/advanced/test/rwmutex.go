package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter 安全的计数器
type SafeCounter struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewSafeCounter 创建新的计数器
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		data: make(map[string]int),
	}
}

// Inc 增加指定 key 的计数（写操作）
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock() // 写锁：独占
	defer c.mu.Unlock()

	c.data[key]++
	fmt.Printf("  [写] %s 增加到 %d\n", key, c.data[key])
}

// Value 获取指定 key 的值（读操作）
func (c *SafeCounter) Value(key string) int {
	c.mu.RLock() // 读锁：允许多个并发读
	defer c.mu.RUnlock()

	return c.data[key]
}

// GetStatistics 获取所有统计信息（批量读操作）
func (c *SafeCounter) GetStatistics() map[string]int {
	c.mu.RLock() // 读锁：保证批量读取的一致性
	defer c.mu.RUnlock()

	// 复制数据，避免外部修改
	result := make(map[string]int)
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

func main() {
	counter := NewSafeCounter()
	var wg sync.WaitGroup

	// 启动 5 个写操作的 goroutine
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				counter.Inc("page_views")
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// 启动 10 个读操作的 goroutine（读操作可以并发）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				val := counter.Value("page_views")
				fmt.Printf("  [读%d] page_views = %d\n", id, val)
				time.Sleep(5 * time.Millisecond)
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 最终读取统计信息（批量读）
	fmt.Println("\n=== 最终统计 ===")
	stats := counter.GetStatistics()
	for k, v := range stats {
		fmt.Printf("%s: %d\n", k, v)
	}
}
