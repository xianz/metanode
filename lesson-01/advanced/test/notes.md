#### select
* 没有default会阻塞，如果ch1和ch2都阻塞，select会一直等待
* 有default，发现阻塞会立即执行

* 关闭channel的处理
* * 从关闭的channel中接收数据，会立即返回数据类型默认值、ok为false
* * 将关闭的channel设置为nil，select会忽略此case
* * 向关闭的channel发送数据，会panic

