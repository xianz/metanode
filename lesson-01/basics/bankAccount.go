package main

/****
* 银行账户系统
 */

import (
	"errors"
	"fmt"
	"strconv"
)

type Account struct {
	Username string
	Balance  int
}

type BankAccount struct {
	Accounts map[string]*Account
	history  []string
}

func (b *BankAccount) AddAccount(username string) {
	b.Accounts[username] = &Account{Username: username, Balance: 0}
}

func (b *BankAccount) Deposit(username string, amount int) (bool, error) {
	if amount < 1 {
		msg := "金额不能小于1"
		b.history = append(b.history, msg)
		return false, errors.New(msg)
	}
	if b.Accounts[username] == nil {
		msg := fmt.Sprintf("[%s] 用户不存在，无法存款", username)
		b.history = append(b.history, msg)
		return false, errors.New(msg)
	}
	b.Accounts[username].Balance += amount
	b.history = append(b.history, "["+username+"]存入"+strconv.Itoa(amount))
	return true, nil
}

func (b *BankAccount) Withdrawal(username string, amount int) (bool, error) {
	user := b.Accounts[username]
	if user.Balance < amount {
		msg := fmt.Sprintf("[%s]账户余额不足，无法取款", username)
		b.history = append(b.history, msg)
		return false, errors.New(msg)
	}
	user.Balance -= amount
	b.history = append(b.history, fmt.Sprintf("[%s]已取出%d，余额：%d", username, amount, user.Balance))
	return true, nil
}

func (b *BankAccount) Info() {
	fmt.Printf("账户中总共有 %d 个用户\n", len(b.Accounts))
}

func main() {
	bank := BankAccount{Accounts: make(map[string]*Account)}

	bank.AddAccount("张三")
	bank.Deposit("张三", 100)
	bank.Deposit("张三", -1)
	bank.Deposit("李四", 150)

	bank.Withdrawal("张三", 80)
	bank.Withdrawal("张三", 30)
	for _, msg := range bank.history {
		fmt.Println(msg)
	}
	bank.Info()

}
