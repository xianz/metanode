package main

import "fmt"

/**
interface 测试
*/

type Database interface {
	Connect() error
	Query(sql string) (interface{}, error)
	Close() error
}

type Mysql struct {
	connection string
}

func (server *Mysql) Connect() error {
	return nil
}

func (server *Mysql) Query(sql string) (interface{}, error) {
	return []string{"row1", "row2"}, nil
}
func (server *Mysql) Close() error {
	return nil
}

func main() {
	mysql := &Mysql{}
	exec(mysql, "select * from talbe1")
}

func exec(db Database, sql string) {
	db.Connect()
	defer db.Close()
	result, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	fmt.Println("Result:", result)
}
