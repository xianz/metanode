package main

import (
	"fmt"
	"strings"
)

type Student struct {
	Name string
}

func main() {

	str := "hello 世界"
	for _, char := range str {
		// fmt.Printf("%s", char)
		fmt.Print(char, " - ")
	}

	fmt.Println()

	src := "hello 世界"
	for _, char := range src {
		fmt.Printf("%c", char)
	}

	m1 := make(map[string]*Student)
	m1["aaa"] = &Student{}
	fmt.Println(m1["bb"])
	return

	fmt.Println("----------")
	var s []Student
	sa := []Student{}
	fmt.Println("student is nil:", s == nil, "len:", len(s), "cap:", cap(s))
	fmt.Println("student slice is nil:", sa == nil, "len:", len(sa), "cap:", cap(sa))
	s = append(s, Student{})
	sa = append(sa, Student{})
	fmt.Println("student is nil:", s == nil, "len:", len(s), "cap:", cap(s))
	fmt.Println("student slice is nil:", sa == nil, "len:", len(sa), "cap:", cap(sa))

	sb := make([]int, 0)
	fmt.Println("sb is nil:", sb == nil, "len:", len(sb), "cap:", cap(sb))
	sb = append(sb, 2)
	fmt.Println("sb is nil:", sb == nil, "len:", len(sb), "cap:", cap(sb))
	for i, value := range sb {
		fmt.Printf("index: %d, value: %v\n", i, value)
	}

	sc := make(map[string]int)
	fmt.Println("sc is nil:", sc == nil, sc, "len:", len(sc))
	sc["a"] = 1
	for key, value := range sc {
		fmt.Printf("key: %s, value: %v\n", key, value)
	}

	f := closure()
	f()
	f()

	add := curryAdd(2)
	fmt.Println("curry add result:", add(3))

	multipartial := curryMultiply(1)(2)(3)
	fmt.Println("curry multiply result:", multipartial)

	logger := createLogger("notify")
	logger("This is a notification message.")
	logger2 := createLogger("error")
	logger2("This is an error message.")

	repeator := repeatstring("#")
	repeator(5)
	repeator(10)

	isInRange := isRange(2, 8)
	fmt.Println("is 4 in range:", isInRange(4))
	fmt.Println("is 1 in range:", isInRange(1))
}

func isRange(min int, max int) func(int) bool {
	return func(value int) bool {
		return value >= min && value <= max
	}
}

func repeatstring(str string) func(int) {
	return func(count int) {
		fmt.Println(strings.Repeat(str, count))
	}
}

func createLogger(prefix string) func(string) {
	return func(message string) {
		fmt.Printf("[%s] %s\n", prefix, message)
	}
}

func curryMultiply(a int) func(int) func(int) int {
	return func(b int) func(int) int {
		return func(c int) int {
			return a + b + c
		}
	}
}

func curryAdd(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

func closure() func() {
	i := 0
	return func() {
		i++
		fmt.Println("i:", i)
	}
}
