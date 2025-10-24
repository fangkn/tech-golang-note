package main

import "fmt"

func add() func(int) int {
	var x int
	return func(y int) int {
		x += y
		return x
	}
}

func main() {
	var f = add()

	fmt.Println(f(10))
	fmt.Println(f(20))
	fmt.Println(f(30))
	fmt.Println("------------")

	f1 := add()
	fmt.Println(f1(40))
	fmt.Println(f1(50))
}
