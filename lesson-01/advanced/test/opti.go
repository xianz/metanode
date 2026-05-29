package main

import (
	"fmt"
	"sync"
)

func workerPool(items []int, numWorkers int) {
	fmt.Println("=== workerPool ===")
	jobs := make(chan int, len(items))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for item := range jobs {
				fmt.Printf("worker %d 处理任务 %d\n", id, item)
			}
		}(i)
	}

	// 发送任务
	for _, item := range items {
		jobs <- item
	}
	close(jobs)
	wg.Wait()
}

func main() {
	workerPool([]int{1, 2, 3, 4, 5, 6, 7, 8}, 3)
}
