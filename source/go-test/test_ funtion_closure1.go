package main

import "fmt"

func cal(base int) (func(int) int, func(int) int) {
	add := func(a int) int {
		base += a
		return base
	}
	sub := func(a int) int {
		base -= a
		return base
	}
	return add, sub
}

func main() {

	add, sub := cal(100)
	r := add(100)
	fmt.Printf("r:%v\n", r)
	r = sub(50)
	fmt.Printf("r:%v\n", r)

	add1, sub1 := cal(100)
	r = add1(1)
	fmt.Printf("r:%v\n", r)
	r = sub1(2)
	fmt.Printf("r:%v\n", r)
}
