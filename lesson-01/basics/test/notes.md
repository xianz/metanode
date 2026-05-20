#### switch语句增加穿透性
case中最后添加一句“fallthrough”，可继续执行下一个case

#### defer 的延迟执行时机
* 函数正常返回前：执行完return语句之后，在函数真正返回之前
* 发生panic时，defer也会执行
* 函数执行到最后一行后，函数自然结束时（没有明确的return语句）

#### panic 和 recover
panic("发生panic") 和 recover() 捕获

#### 函数返回
* 多返回值：func f1() (int, error){ return 0, errors.New("...")}
* 命名返回：func f2() (a, b int) { a,b= 2, 3; return}