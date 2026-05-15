package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	////多返回值
	// fmt.Println(divide(3, 2))
	//// 命名值返回
	// a, b := f2()
	// fmt.Println(a, b)
	//// 可变参数
	// f3(1, 2, 3, 4, 5)
	// f3() // 没有参数也是合法的
	//// 闭包
	// counter := closureExample()
	// counter() // 当前计数: 1
	// counter() // 当前计数: 2
	// counter() // 当前计数: 3
	//// 柯里化
	// add1 := makeAddr(2)
	// add1(2) //
	// add2 := makeAddr(1)
	// add2(1)
	// add1(2)
	//// 字符串重复函数
	repeatSplit := repeatString(22)
	repeatSplit("-")
}

// 多返回值
func divide(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

// 命名返回值
func f2() (a, b int) {
	a, b = 1, 2
	return
}

// 可变参数
func f3(ns ...int) {
	fmt.Printf("\n传入的参数个数: %d，分别是：", len(ns))
	for _, v := range ns {
		fmt.Printf("%d ", v)
	}
}

// 闭包
func closureExample() func() {
	count := 0
	return func() {
		count++
		fmt.Println("当前计数:", count)
	}
}

// 柯里化(闭包的一个应用)
func makeAddr(defaultValue int) func(n int) int {
	count := defaultValue
	return func(n int) int {
		count += n
		fmt.Println("count的值是：", count)
		return count
	}
}

// 另一个柯里化示例：字符串重复函数
func repeatString(n int) func(s string) {
	return func(s string) {
		fmt.Println(strings.Repeat(s, n))
	}
}
