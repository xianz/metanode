package main

import (
	"fmt"
	"log"
)

var dbName = "users"

func main() {
	// 得到数据库
	db, err := SetupDB(dbName)
	if err != nil {
		log.Fatalln(err)
	}

	// 更新表结构
	if err := UpdateDB(db); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("添加字段完成")

	// // 新增用户
	// if _, err := CreateUser(db, "张三", "zs@qq.com"); err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("新增用户完成")

	// 模糊查询
	var emailPattern = "%2@%"
	var size = 3
	var page = 1
	users, err := SearchUsersByEmail(db, emailPattern, page, size)
	if err != nil {
		log.Fatalln(err)
	}
	for _, u := range users {
		fmt.Println(u)
	}

	// 批量更新状态
	ids := []uint{2, 3, 4}
	err = UpdateUserStatus(db, ids, "pending")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("批量更新完成！！")

	// 删除过期用户
	if err := DeleteInactiveUsers(db); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("删除完成，影响行数：", db.RowsAffected)

	// 分页 scopes
	var userResult []User
	db.Scopes(ScopeAge(18, 30), ScopePage(1, 3)).Find(&userResult)
	fmt.Println(userResult)
}
