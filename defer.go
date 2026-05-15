package main

import "fmt"

func main() {
	// fmt.Println(deferClosure())
	// panicTest()
	panicAndRecover()
}

func deferClosure() int {
	i := 0
	defer func() {
		fmt.Println("defer闭包，i=", i)
	}()
	i++
	return i // 不能
}

func panicTest() {
	defer func() {
		fmt.Println("panic不耽误执行defer")
	}()
	var p *int = nil
	fmt.Println(*p) // 这里会引发panic
}

func panicAndRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover捕获到panic:", r)
		}
	}()
	fmt.Println("函数体正常执行")
	panic("发生了一个自定义错误") // 这里会引发panic
}
