package main

import "fmt"

var i int = initVer()

func init() {
	fmt.Println("0 : init......")
}

func initVer() int {
	fmt.Println("0 : initVer.....")
	return 100
}

func main() {
	fmt.Println("1 : main......")
}
