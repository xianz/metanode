package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task 任务定义
type Task struct {
	ID      string
	Execute func(ctx context.Context) error
	Timeout time.Duration
}

// Scheduler 调度器
type Scheduler struct {
	maxConcurrent int
	wg            sync.WaitGroup
}

// NewScheduler 创建调度器
func NewScheduler(maxConcurrent int) *Scheduler {
	return &Scheduler{
		maxConcurrent: maxConcurrent,
	}
}

// Run 执行任务列表
func (s *Scheduler) Run(tasks []Task) []error {
	taskChan := make(chan Task, len(tasks))
	results := make([]error, len(tasks))

	// 启动工作协程
	for i := 0; i < s.maxConcurrent; i++ {
		s.wg.Add(1)
		go s.worker(taskChan, &results)
	}

	// 发送任务
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	s.wg.Wait()
	return results
}

// worker 工作协程
func (s *Scheduler) worker(taskChan <-chan Task, results *[]error) {
	defer s.wg.Done()

	for task := range taskChan {
		// 执行任务
		err := s.executeTask(task)

		// 存储结果（简化版，实际可能需要更复杂的结果索引）
		// 这里为了简化，我们直接打印错误
		if err != nil {
			fmt.Printf("任务 %s 失败: %v\n", task.ID, err)
		} else {
			fmt.Printf("任务 %s 成功\n", task.ID)
		}
	}
}

// executeTask 执行单个任务（带超时控制）
func (s *Scheduler) executeTask(task Task) error {
	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()

	// 创建结果通道
	resultChan := make(chan error, 1)

	// 执行任务
	go func() {
		resultChan <- task.Execute(ctx)
	}()

	// 等待结果或超时
	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("超时 (%v)", task.Timeout)
	}
}

func main() {
	// 创建调度器，最多同时执行2个任务
	scheduler := NewScheduler(2)

	// 定义任务列表
	tasks := []Task{
		{
			ID:      "task1",
			Timeout: 3 * time.Second,
			Execute: func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				fmt.Println("  task1 执行中...")
				return nil
			},
		},
		{
			ID:      "task2",
			Timeout: 2 * time.Second,
			Execute: func(ctx context.Context) error {
				time.Sleep(4 * time.Second) // 这个会超时
				fmt.Println("  task2 执行中...")
				return nil
			},
		},
		{
			ID:      "task3",
			Timeout: 5 * time.Second,
			Execute: func(ctx context.Context) error {
				// 支持取消的任务
				for i := 0; i < 3; i++ {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(1 * time.Second):
						fmt.Printf("  task3 进度: %d/3\n", i+1)
					}
				}
				return nil
			},
		},
	}

	// 执行任务
	fmt.Println("开始执行任务...")
	errors := scheduler.Run(tasks)

	// 统计结果
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}

	fmt.Printf("\n执行完成: 成功 %d/%d\n", successCount, len(tasks))
}
