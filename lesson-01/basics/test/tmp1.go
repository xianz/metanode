package main

import "fmt"

type T struct{}
type PT *T

// 方法集规则：
// - T 的方法集：接收者为 T 的方法
// - *T 的方法集：接收者为 T 和 *T 的方法
func (T) ValMethod()  {} // 值接收者
func (*T) PtrMethod() {} // 指针接收者

// 接口实现的关键差异：
var i interface{ ValMethod() }
var i2 interface{ PtrMethod() }

func main() {
	var t T
	var pt *T
	fmt.Printf("pt:%#v\n", pt)

	t.ValMethod()
	t.PtrMethod()  // 编译器自动取址
	pt.ValMethod() // 编译器自动解引用
	pt.PtrMethod()

	i = t // T 实现了接口

	i = pt // *T 也实现了接口

	//i2 = t // T 没有实现！只有 *T 实现了
	i2 = pt

}
