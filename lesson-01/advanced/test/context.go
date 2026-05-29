package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// contextCancelDemo()
	// contextDeadLineDemo()
	// contextValueDemo()
	// contextChildcancelDemo()
	multiWorkDemo()
}

// 手动取消
func contextCancelDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("完成")
				return
			default:
				fmt.Println("正在工作...")
				time.Sleep(time.Second)
			}
		}
	}()
	time.Sleep(time.Second * 3)
	fmt.Println("取消 enter")
	cancel()
	fmt.Scanln()
}

// 绝对时间
func contextDeadLineDemo() {
	fmt.Println("=== contextDeadLine ===")
	deadLineTime := time.Now().Add(time.Second * 1)
	ctx, _ := context.WithDeadline(context.Background(), deadLineTime)
	defer func() {
		fmt.Println("defer了")
		// cancel()  // 不进行收回也不会报错
	}()

	select {
	case <-time.After(time.Second * 2):
		fmt.Println("到达2秒")
	case <-ctx.Done():
		fmt.Println("到达DeadLine")
	}
	fmt.Scanln()
}

// 传值
func contextValueDemo() {
	ctx := context.WithValue(context.Background(), "a", "aastr")
	ctx = context.WithValue(ctx, "b", "bbstr")
	processRequest(ctx)
}
func processRequest(ctx context.Context) {
	if value := ctx.Value("a"); value != nil {
		fmt.Println(value)
	} else {
		fmt.Println("没取到 a ")
	}
	if value := ctx.Value("b"); value != nil {
		fmt.Println(value)
	} else {
		fmt.Println("没取到 a ")
	}
}

// 级联取消
func contextChildcancelDemo() {
	ctx, parentCancel := context.WithCancel(context.Background())
	childCtx1, child1Cancel := context.WithCancel(ctx)
	defer child1Cancel()
	childCtx2, child2Cancel := context.WithCancel(ctx)
	defer child2Cancel()
	go worker(childCtx1, "child 1")
	go worker(childCtx2, "child 2")

	time.Sleep(time.Second * 3)
	parentCancel()
	fmt.Scanln()
}
func worker(ctx context.Context, tag string) {
	for i := 0; i < 10; i++ { // 模拟持续工作
		select {
		case <-ctx.Done():
			fmt.Println(tag, "已手动取消")
			return
		default:
			fmt.Println(tag, "持续工作中...")
			time.Sleep(time.Second)
		}
	}
}

// 取消多个goroutine
func multiWorkDemo() {
	fmt.Println("=== 多worker ===")
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 3; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("worker:", id, "退出\n")
					return
				default:
					fmt.Println("worker", id, "工作中...")
					time.Sleep(time.Second)
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 2)
	cancel()
	fmt.Scanln()
}
