package main

import "fmt"

func doSomething(val interface{}) {
	fmt.Printf("值：%v，类型：%T\n", val, val)
	// 断言：普通类型断言（返回一个值+是否成功）
	// value, ok := string(val)

	// 断言：类型开关，判断多种类型
	switch data := val.(type) {
	case int:
		fmt.Println("int类型:", data)
	case string:
		fmt.Println("string类型:", data)
	case struct{}:
		fmt.Println("struct{}类型：", data)
	case []int:
		fmt.Println("int数组类型：", data)
	case interface{}: // 只能放到最后面，否则会优先匹配
		fmt.Println("interface类型:", data)
	default:
		fmt.Println("位置类型:", data)

	}
}

func main() {

	payload := map[string]interface{}{
		"id":    1001,
		"name":  "golang",
		"extra": []string{"aaaa", "bbb"},
	}

	if id, ok := payload["id"].(int); ok {
		fmt.Println("ID:", id)
	}
	if tags, ok := payload["extra"].([]string); ok {
		fmt.Println("tags:", tags)
	} else {
		fmt.Println("extra不是期望的[]string")
	}

	var a interface{}
	fmt.Println("=== 类型switch type ===")
	// a = 3
	// i, ok := a.(type)
	// fmt.Println(i, ok)

	// var a interface{}
	// doSomething(a)
	// a = 3
	// doSomething(a)
	// a = "abc"
	// doSomething(a)

	a = Dog{Name: "gogo"}
	doSomething(a)

	a = []int{3, 2, 1}
	doSomething(a)
}

// type Animal interface{ sleep() }
type Dog struct{ Name string }

// func (d *Dog) sleep() {
// 	fmt.Println("dog sleep")
// }
