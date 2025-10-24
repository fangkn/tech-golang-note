package main

import (
	"bytes"
	"fmt"
)

func myfunc(argv ...int) {
	for _, arg := range argv {
		fmt.Println(arg)
	}
}

func MyPrintf(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Println(arg, "is an int value.")
		case string:
			fmt.Println(arg, "is a string value.")
		case int64:
			fmt.Println(arg, "is an int64 value.")
		default:
			fmt.Println(arg, "is an unknown type.")
		}
	}
}

func joinStrings(slist ...string) string {
	// 定义一个字节缓冲, 快速地连接字符串
	var b bytes.Buffer
	// 遍历可变参数列表slist, 类型为[]string
	for _, s := range slist {
		// 将遍历出的字符串连续写入字节数组
		b.WriteString(s)
	}
	// 将连接好的字节数组转换为字符串并输出
	return b.String()
}

// 实际打印的函数
func rawPrint(rawList ...interface{}) {
	// 遍历可变参数切片
	for _, a := range rawList {
		// 打印参数
		fmt.Println(a)
	}
}

// 打印函数封装
func print(slist ...interface{}) {

	rawPrint(slist...) // 将slist可变参数切片完整传递给下一个函数
}

func main() {
	// myfunc(1, 3, 4, 5)
	// myfunc(10, 30, 40, 50)
	// myfunc(100, 300, 400, 500)

	// var v1 int = 1
	// var v2 int64 = 234
	// var v3 string = "hello"
	// var v4 float32 = 1.234
	// MyPrintf(v1, v2, v3, v4)

	// 输入3个字符串, 将它们连成一个字符串
	//fmt.Println(joinStrings("pig ", "and", " rat"))
	//fmt.Println(joinStrings("hammer", " mom", " and", " hawk"))

	print(1, 2, 3)

}
