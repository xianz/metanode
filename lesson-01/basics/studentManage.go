package main

/**
学生管理：添加、修改、删除
*/

import (
	"fmt"
)

type Student struct {
	Name  string
	Score int
}

func (s *Student) GetInfo() {
	println("姓名：", s.Name, "分数：", s.Score)
}

type StudentManager struct {
	students map[string]*Student // 学生对象必须是指针类型，否则修改分数是只会修改副本的分数
}

func (sm *StudentManager) Init() {
	sm.students = make(map[string]*Student)
}

func (sm *StudentManager) AddStudent(s Student) {
	_, exists := sm.students[s.Name]
	if exists {
		fmt.Println("学生已存在")
		return
	}
	sm.students[s.Name] = &s
	fmt.Println("添加学生成功")
}

func (sm *StudentManager) DeleteStudent(name string) {
	delete(sm.students, name)
	fmt.Println("删除学生成功")
}

func (sm *StudentManager) ModifyScore(name string, score int) {
	student, exists := sm.students[name]
	if exists == false {
		fmt.Println("学生不存在，无法修改分数")
		return
	}
	student.Score = score
	fmt.Println("修改学生成绩成功")
}

func (sm *StudentManager) StudentList() {
	fmt.Println("学生列表：")
	for _, student := range sm.students {
		student.GetInfo()
	}
}

func main() {
	sm := StudentManager{}
	sm.Init()
	// 添加学生
	fmt.Println("--------- 添加学生")
	sm.AddStudent(Student{Name: "张三", Score: 80})
	sm.AddStudent(Student{Name: "张三", Score: 90}) // 测试添加重复学生
	sm.AddStudent(Student{Name: "里斯", Score: 55})
	// 显示学生列表
	sm.StudentList()

	fmt.Println("--------- 修改学生分数")
	sm.ModifyScore("张三", 90) // 修改张三的分数为90
	sm.ModifyScore("李四", 70) // 测试修改不存在的学生分数
	sm.StudentList()

	fmt.Println("--------- 删除学生")
	sm.DeleteStudent("里斯")
	sm.StudentList()
}
