package main

import "fmt"

func main() {
	var i int = 42 // 声明一个整数变量i，初始值为42
	// var p *int = &i // 声明一个指向整数的指针，初始值为nil
	changeValue(&i)
	fmt.Println(i)
	p := &i
	fmt.Println(p)

	////
	if n := 3; n > 0 {
		fmt.Println("n is greater than 0")
	}

	//// for循环
	scores := []int{90, 80, 70}
	for index, score := range scores {
		fmt.Printf("Index: %d, Score: %d\n", index, score)
	}

	//// switch语句
	age := 18
	switch { // 不穿透性
	case age == 18:
		fmt.Println("You are 18 years old.")
	case age == 17:
		fmt.Println("You are 17 years old.")
		fallthrough // 穿透性，继续执行下一个case
	case age <= 16:
		fmt.Println("You are 16 years old.")
	default:
		fmt.Println("必须小于 18 岁")
	}

	fmt.Println("-------------------")
	src := "hello 世界" // 字符串是一个字节序列，for range会将每个Unicode字符 作为rune类型的值返回。rune是Go语言中的一个类型，表示一个Unicode字符，可以存储一个字符的Unicode码点。对于ASCII字符，rune的值与字符的ASCII码相同；对于非ASCII字符，rune的值是该字符的Unicode码点。
	for _, char := range src {
		fmt.Printf("%c ", char)
	}

	m1 := map[string]int{"Alice": 30, "Bob": 25}
	for name, age := range m1 {
		fmt.Printf("name: %s, age:%d\n", name, age)
	}
}

func changeValue(p *int) {
	*p = 100 // 通过指针修改i的值
}
