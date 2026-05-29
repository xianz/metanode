package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func gmpDemo() {
	fmt.Println("=== GMP调度示例===")
	fmt.Println("逻辑CPU数量：", runtime.NumCPU())
	// 如果 n>0，设置可使用的数量，返回旧值
	// 如果n==0，不修改GOMAXPROCS，只查询并返回
	prev := runtime.GOMAXPROCS(0)
	fmt.Println("当前可用GoProcs：", prev)
	runtime.GOMAXPROCS(2) // 设置为用2个

	var wg sync.WaitGroup
	start := time.Now()
	taskCount := 8
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		go func(id int) {
			defer wg.Done()
			// 模拟CPU密集型工作
			sum := 0
			for j := 0; j < 50_000_0000; j++ {
				sum += j % (id + 1)
			}
			fmt.Printf("Go Routine %d 完成，结果：%d\n", id, sum)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\n所有任务完成，耗时：%s\n", elapsed)
	fmt.Println("提示：运行时可配合命令 `GODEBUG=schedtrace=1000,scheddetail=1 go run 06-gmp.go` 观察调度日志。")
	// 恢复原始的 GOMAXPROCS，避免影响其他程序
	runtime.GOMAXPROCS(prev)
}

func main() {
	gmpDemo()
}
