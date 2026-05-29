package main

import (
	"errors"
	"fmt"
	"time"
)

/**
接口练习
- 实现一个支付系统
- 支持多种支付方式（支付宝、微信、银行卡）
- 使用接口实现多态
*/
// 支付接口
type Payment interface {
	Pay(amount float64) error
	Unpay(amount float64) error
}

// 支付方法一：支付宝
type AliPay struct {
	UserName string
	Status   string
	PayTime  time.Time
}

func (a *AliPay) Pay(amount float64) error {
	if amount <= 0 {
		return errors.New("付款金额不能<=0")
	}
	fmt.Printf("【支付宝】账户[%s]，付款完成\n", a.UserName)
	a.Status = "已付款"
	a.PayTime = time.Now()
	return nil
}
func (a *AliPay) Unpay(amount float64) error {
	if amount <= 0 {
		return errors.New("退款金额不正确")
	}
	fmt.Println("【支付宝】退款完成")
	a.Status = "已退款"
	return nil
}

// 支付方法二：微信
type WechatPay struct {
	openId  string
	Status  string
	payTime time.Time
}

func (w *WechatPay) Pay(amount float64) error {
	if amount <= 0 {
		return errors.New("付款金额不能小于或等于0")
	}
	fmt.Printf("【微信支付】用户[%s], 付款完成\n", w.openId)
	w.Status = "已支付"
	w.payTime = time.Now()
	return nil
}
func (w *WechatPay) Unpay(amount float64) error {
	if amount <= 0 {
		return errors.New("退款金额不正确\n")
	}
	fmt.Printf("【微信支付】退款完成")
	w.Status = "已退款"
	return nil
}

// 支付方式三：银行卡
type CardPay struct {
	CardName   string
	CardNumber string
	Status     string
	PayTime    time.Time
}

func (c *CardPay) Pay(amount float64) error {
	if amount <= 0 {
		return errors.New("支付金额不对！\n")
	}
	fmt.Printf("【银行卡】卡号[%s]，支付完成\n", c.CardNumber)
	c.Status = "已支付"
	return nil
}
func (c *CardPay) Unpay(amount float64) error {
	fmt.Println("【银行卡】退款完成")
	c.Status = "已退款"
	return nil
}

// 支付订单
type Order struct {
	OrderId     uint
	Amount      float64
	payment     Payment
	CrementTime time.Time
}

func (o *Order) SetPayment(payment Payment) {
	o.payment = payment
}

func (o *Order) Checkout() {
	fmt.Printf("==== 订单 %d ====\n", o.OrderId)
	err := o.payment.Pay(o.Amount)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	alipay := &AliPay{UserName: "zhangsan@qq.com", Status: "待付款"}
	// 支付宝
	order1 := &Order{OrderId: 1010110, Amount: 3.05}
	order1.SetPayment(alipay)
	order1.Checkout()
	order1.payment.Unpay(2)
	// 微信
	wxpay := &WechatPay{openId: "op_id123", Status: "待付款"}
	order2 := &Order{OrderId: 1010111, Amount: 0}
	order2.SetPayment(wxpay)
	order2.Checkout()
	order2.Amount = 2.4
	order2.Checkout()
	//银行卡
	order3 := &Order{OrderId: 1010111, Amount: 0.2}
	cardPay := &CardPay{CardName: "张三", CardNumber: "2424551545545", Status: "待付款"}
	order3.SetPayment(cardPay)
	order3.Checkout()

}
