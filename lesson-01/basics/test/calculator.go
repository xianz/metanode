package main

import (
	"errors"
	"fmt"
)

type Calculator struct {
	logs []string
}

// 加法：支持多个数字相加
func (c *Calculator) Add(numbers ...float64) (result float64, msg error) {
	defer func() {
		c.logs = append(c.logs, fmt.Sprintf("Add: %v = %v，%s", numbers, result, msg))
	}()
	if len(numbers) < 2 {
		msg = errors.New("至少需要两个数字进行加法运算")
		return
	}
	for _, num := range numbers {
		result += num
	}
	return
}

// 减法
func (c *Calculator) Subtract(a, b float64) (result float64, msg error) {
	c.logs = append(c.logs, fmt.Sprintf("Subtract: %v - %v = %v", a, b, a-b))
	return a - b, nil
}

// 乘法
func (c *Calculator) Multiply(a, b float64) (result float64, msg error) {
	c.logs = append(c.logs, fmt.Sprintf("Multiply: %v * %v = %v", a, b, a*b))
	return a * b, nil
}

// 除法：需要处理除数为零的情况
func (c *Calculator) Divide(a, b float64) (result float64, msg error) {
	defer func() {
		c.logs = append(c.logs, fmt.Sprintf("Divide: %v / %v = %v，%s", a, b, result, msg))
	}()
	if b == 0 {
		msg = errors.New("除数不能为零")
		return
	}
	return a / b, nil
}

// 平均值：支持多个数字求平均值，并处理输入为空的情况
func (c *Calculator) Average(numbers ...float64) (result float64, msg error) {
	defer func() {
		c.logs = append(c.logs, fmt.Sprintf("Average: %v = %v，%s", numbers, result, msg))
	}()
	if len(numbers) < 1 {
		msg = errors.New("至少需要一个数字进行平均值计算")
		return
	}
	for _, num := range numbers {
		result += num
	}
	result = result / float64(len(numbers))
	return
}

func main() {
	// 创建计算器实例
	calc := &Calculator{}
	//// 加法
	result, _ := calc.Add(1, 2.3)
	fmt.Println("1 + 2.3 =", result)

	//// 减法
	result, _ = calc.Subtract(5, 3)
	fmt.Println("5 - 3 = ", result)

	//// 乘法
	result, _ = calc.Multiply(33, 2)
	fmt.Println("33 * 2 = ", result)

	//// 正常除法
	result, _ = calc.Divide(10, 3)
	fmt.Println("10 / 3 = ", result)

	//// 除数为零
	result, err := calc.Divide(10, 0)
	if err != nil {
		fmt.Println("10 / 0 ", err)
	}

	//// 求平均值：正常情况
	result, _ = calc.Average(1, 2, 3, 4)
	fmt.Println("平均值(1, 2, 3, 4) = ", result)

	//// 求平均值：错误情况
	result, err = calc.Average()
	if err != nil {
		fmt.Println("Average ", err)
	}

}
