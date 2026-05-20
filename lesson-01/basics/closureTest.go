package main

/**
* 3 个使用闭包的实际例子
 */

import (
	"fmt"
	"strings"
)

func main() {

	// 例子1：生成分隔符
	repeator := stringRepeat("-")
	repeator(33)

	// 例子2：惰性斐波那契数列
	f := fibonacci()
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	repeator(33)

	// 例子3：区间过滤器
	filter := greaterThan(2, 6)
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	for _, n := range numbers {
		if filter(n) {
			fmt.Print(n, " ")
		}
	}
}

// 重复字符串
func stringRepeat(str string) func(int) {
	return func(count int) {
		fmt.Println(strings.Repeat(str, count))
	}
}

// 惰性斐波那契数列
func fibonacci() func() int {
	x, y := 0, 1

	return func() int {
		// // 方式一
		// defer func() {
		// 	x, y = y, x+y
		// }()
		// return x + y
		// 方式二
		result := x
		x, y = y, x+y
		return result
	}
}

// 判断指定数字是否包含在n，m区间
func greaterThan(n, m int) func(int) bool {
	return func(i int) bool {
		return n < i && i < m
	}
}
