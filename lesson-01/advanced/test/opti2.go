package main

import "fmt"

func workerPool(items []int, numWorkers int) {
	jobs := make(chan int, len(items))

	// 启动worker
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for item := range jobs {
				fmt.Printf("Worker %d 处理任务 %d\n", workerID, item)
			}
		}(w)
	}

	// 发送任务
	for _, item := range items {
		jobs <- item
	}
	close(jobs)
}

func main() {
	// 准备5个任务
	items := []int{1, 2, 3, 4, 5}

	// 启动3个worker的工作池
	workerPool(items, 3)

	// 等待goroutine执行完毕
	fmt.Scanln()
}
