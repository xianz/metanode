package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (s *SafeCounter) Increment(m int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	fmt.Printf("goroutine %d incremented count %d\n", m, s.count)
}

func (s *SafeCounter) GetCount() int {
	return s.count
}

func mutexDemo() {
	fmt.Println("==== Mutex test ====")
	sc := &SafeCounter{}

	var sw sync.WaitGroup
	for i := 0; i < 30; i++ {
		sw.Add(1)
		go func(i int) {
			defer sw.Done()
			sc.Increment(i)
		}(i) // 避免共享外部i
	}
	sw.Wait()
	fmt.Println("最终为：", sc.GetCount())
}

type rwSafeCounter struct {
	mu    sync.RWMutex
	count int
}

func (s *rwSafeCounter) RWMutexDemo() {
	fmt.Println(" ===RWMutex test===")
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 10)
		go func(i int) {
			s.mu.Lock()
			defer func() {
				s.mu.Unlock()
				wg.Done()
			}()
			s.count++
			fmt.Printf("[写]%d 增加到: %v\n", i, s.count)
		}(i)
	}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 9)
		go func(i int) {
			s.mu.RLock()
			defer func() {
				s.mu.RUnlock()
				wg.Done()
			}()
			fmt.Printf("[读]%d count: %v\n", i, s.count)
		}(i)
	}
	wg.Wait()
}

func contextTimeoutDemo() {
	ct, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ch := make(chan string)
	go func() {
		time.Sleep(time.Second * 3)
		ch <- "body"
	}()
	select {
	case data := <-ch:
		fmt.Println(data)
	case <-ct.Done():
		fmt.Println("超时了：", ct.Err())
		// case <-time.After(time.Second * 4):  // 超时方法二
		// 	fmt.Println("超时了")
	}
}

func contextDeadLineDemo() {
	fmt.Println("=== deadLine context ===")
	ct, cancel := context.WithDeadline(context.Background(),
		time.Now().Add(time.Second*3))
	defer cancel()

	// 检查截止时间，做相应业务逻辑
	go func() {
		time.Sleep(time.Millisecond * 500)
		date, ok := ct.Deadline() // 返回值：截止时间，是否设置了截至时间
		if ok {
			fmt.Printf("截至时间：%v，剩余：%v\n", date, time.Until(date))
		}
	}()

	select {
	case <-time.After(time.Second * 4):
		fmt.Println("模拟工作完成")
	case <-ct.Done():
		fmt.Println("超过截至时间了", ct.Err())
	}
	fmt.Println()
}

func contextValueDemo() {
	fmt.Println("=== content value ===")
	// 传递
	ct := context.WithValue(context.Background(), "reqID", "req123")
	ct = context.WithValue(ct, "userID", "user123")

	// 取值
	val := ct.Value("reqID")
	if val != nil {
		fmt.Println(val)
	}
}

// 超时后的错误处理
func contextErrorHandlingDemo() {
	fmt.Println("=== context错误处理 ===")
	ct, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	time.Sleep(time.Second * 1)

	// 可包装成func
	err := ct.Err()
	switch err {
	case context.Canceled: // 手动调用cancel导致的取消
		fmt.Println("被用户取消")
	case context.DeadlineExceeded:
		fmt.Println("等待超时")
	case nil:
		fmt.Println("没有错误，context仍然有效")
	default: // 其他未知错误，理论上不应该出现（可当作用户检查状态判断使用）
		fmt.Println("未知错误：", err)
	}
}

func main() {
	// mutexDemo()

	// rws := &rwSafeCounter{}
	// rws.RWMutexDemo()

	// contextTimeoutDemo()
	// contextDeadLineDemo()
	// contextValueDemo()
	contextErrorHandlingDemo()

}
