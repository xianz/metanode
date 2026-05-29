package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id      int
	f       func() error
	timeout time.Duration // 任务最多允许执行多久
}

type Scheduler struct {
	MaxConcurrency int // 开启几个协程
	wg             sync.WaitGroup
}

func (s *Scheduler) Run(tasks []Task) {
	taskQueue := make(chan Task, len(tasks))
	// 发送任务
	for _, task := range tasks {
		taskQueue <- task
	}
	close(taskQueue)

	for i := 0; i < s.MaxConcurrency; i++ { // 启动几个goroutine执行，让每个goroutine自己去抢任务做
		s.wg.Add(1)
		go s.worker(taskQueue)
	}
	s.wg.Wait()
}

func (s *Scheduler) worker(taskQueue <-chan Task) {
	defer s.wg.Done()
	for task := range taskQueue { // 每个goroutine都在抢通道内的任务执行
		err := execTask(task)
		if err != nil {
			fmt.Println("执行失败:", err)
		} else {
			fmt.Println("执行成功")
		}
	}
}

func execTask(task Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), task.timeout)
	resultChan := make(chan error)
	defer cancel()
	defer close(resultChan)
	// 执行用户函数 和 等待结果或超时
	go func() {
		resultChan <- task.f()
	}()
	select {
	case rs := <-resultChan:
		return rs // 不一定是错误
	case <-ctx.Done():
		return fmt.Errorf("任务[%d]超时(%v)", task.id, task.timeout)
	}
}

func main() {
	tasks := []Task{
		{id: 1, f: func() error {
			fmt.Println("任务1执行..")
			time.Sleep(time.Second * 4)
			return nil
		}, timeout: time.Second * 3,
		},
		{id: 2, f: func() error {
			fmt.Println("任务2执行..")
			return nil
		}, timeout: time.Second * 2},
		{id: 3, f: func() error {
			fmt.Println("执行任务3...")
			time.Sleep(time.Second)
			return fmt.Errorf("小任务不跑了，停了")
		}, timeout: time.Second * 8},
	}

	sc := &Scheduler{MaxConcurrency: 3}
	sc.Run(tasks)

	// fmt.Scanln()
}
