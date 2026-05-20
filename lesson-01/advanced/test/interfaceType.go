package main

import "fmt"

type Animal interface {
	sleep()
}

type Cat struct {
	Name string
}

func (a *Cat) sleep() {}

type Dog struct {
	Name string
}

func (d *Dog) sleep() {}

func PrintAnimalInfo(animal Animal) {
	fmt.Println(animal)
}

func main() {
	animals := []Animal{&Cat{Name: "小花"}, &Dog{Name: "小黄"}}
	for _, a := range animals {
		PrintAnimalInfo(a)
	}
}
