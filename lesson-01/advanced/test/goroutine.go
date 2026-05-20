package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 等待协程执行完在走主进程（官方推荐的做法） https://chat.deepseek.com/a/chat/s/f902cbc2-82fa-4706-a187-a32c053e1b86
func waitGroupDemo() {
	fmt.Println("====waitGroup 示例 ====")
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("task %d startd \n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("all tasks completed\n")
}

// 无缓冲chan
func channelDemo() {
	ch := make(chan int)
	go func() {
		ch <- 1
		fmt.Println("go routine")
	}()
	// value := <-ch
	fmt.Println(<-ch)
}

// 有缓冲channel
func bufferedChannelDemo() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4

	fmt.Println(<-ch) // 未关闭ch前无法使用 for...range 遍历
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

// 关闭后仍然可读取
func closeChannelDemo() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	for value := range ch {
		fmt.Println(value)
	}
}

// 模拟生产消费测试有缓冲channel
func bufferedChannelDemo2() {
	fmt.Println("==== 生产、消费 ====")
	ch := make(chan int, 3)
	fmt.Println("ready...")
	time.Sleep(time.Second * 1)
	// 写
	go func() {
		for i := 0; i < 30; i++ {
			fmt.Println("写：", i)
			ch <- i
		}
	}()
	fmt.Println("~~")

	// 读1
	go func() {
		for i := 0; i < 15; i++ {
			fmt.Println("读1：", <-ch)
			time.Sleep(time.Second * 2)
		}
	}()
	// 读2
	go func() {
		for i := 0; i < 15; i++ {
			fmt.Println("读2：", <-ch)
			time.Sleep(time.Second * 2)
		}
	}()
	fmt.Println("~~~~")
	time.Sleep(time.Second * 30)
}

// 多路复用，哪个先有数据就处理哪个
func selectDemo() {
	fmt.Println("==== 多路复用 ====")
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		time.Sleep(time.Second * 1)
		ch1 <- 1
	}()
	go func() {
		time.Sleep(time.Second * 1)
		ch2 <- 2
	}()
	go func() {
		time.Sleep(time.Second * 1)
		ch3 <- 3
	}()

	select {
	case i := <-ch1:
		fmt.Println("ch1:", i)
	case i := <-ch2:
		fmt.Println("ch2:", i)
	case i := <-ch3:
		fmt.Println("ch3:", i)
	}
}

// 阻塞超时
func timeoutDemo() {
	fmt.Println("==== 阻塞超时 ====")
	ch := make(chan int)

	go func() {
		time.Sleep(time.Second * 1)
		ch <- 1
	}()

	select {
	case data := <-ch:
		fmt.Println("data:", data)
	case <-time.After(time.Second * 2):
		fmt.Println("超时了")
	}
}

// 非阻塞（加default）
func noBlockingDemo() {
	fmt.Println("==== 非阻塞 ====")
	ch := make(chan int, 1)

	select {
	case ch <- 1:
		fmt.Println("发送成功")
	default:
		fmt.Println("channel已满，发送失败")
	}

	select {
	case data := <-ch:
		fmt.Println("data：", data)
	default:
		fmt.Println("没有值可读取")
	}
}

// 持续监听多个channel，直到关闭
func loopSelect() {
	fmt.Println("==== 监听多个channel ====")
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for i := 0; i < 4; i++ {
			ch1 <- "ch1：" + strconv.Itoa(i)
		}
		close(ch1)
	}()
	go func() {
		for i := 0; i < 6; i++ {
			ch2 <- "ch2：" + strconv.Itoa(i)
		}
		close(ch2)
	}()

	for {
		time.Sleep(time.Millisecond * 100)
		select {
		case val, ok := <-ch1: // 关闭的channel设为nil，select会忽略这个case（是的）
			fmt.Print("ch1 的 case：")
			if !ok {
				ch1 = nil
				fmt.Println("设置了 ch1 为 nil")
				continue
			}
			fmt.Println(val)

		case val, ok := <-ch2:
			fmt.Print("ch2 的 case：")
			if !ok {
				fmt.Println("设置了 ch2 为 nil")
				ch2 = nil
				continue
			}
			fmt.Println(val)
		default:
			// fmt.Println("default；")
			if ch1 == nil && ch2 == nil {
				fmt.Println("所有channel关闭")
				return
			}
		}
	}
}

// 使用退出信号（不适用close(ch)）来关闭
func quitChannel() {
	fmt.Println("==== quitChannel ====")
	ch := make(chan int)
	quit := make(chan bool)

	go func() {
		defer fmt.Println("break ed")
		for {
			select {
			case data := <-ch:
				fmt.Println("data：", data)
				time.Sleep(time.Millisecond * 500)
			case <-quit:
				fmt.Println("收到退出信号")
				// break
				return
			}
		}
	}()

	for i := 0; i < 5; i++ {
		ch <- i
	}
	quit <- true
	fmt.Println("done")

}

func main() {
	// channelDemo()
	// bufferedChannelDemo()
	// bufferedChannelDemo2()
	// waitGroupDemo()
	// closeChannelDemo()
	// selectDemo()
	// timeoutDemo()
	// noBlockingDemo()
	// loopSelect()
	quitChannel()
}
