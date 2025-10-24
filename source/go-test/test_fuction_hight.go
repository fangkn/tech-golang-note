package main

import "fmt"

func saySomething(text string) {
	fmt.Printf("Hello, %s", text)
}

func test(name string, f func(string)) {
	f(name)
}

func add(a int, b int) int {
	return a + b
}

func sub(a int, b int) int {
	return a - b
}

func cal(opt string) func(int, int) int {
	switch opt {
	case "+":
		return add
	case "-":
		return sub
	default:
		return nil
	}
}

func main() {
	// test("carl", saySomething)
	f := cal("+")
	r := f(1, 3)
	fmt.Printf("r:%v", r)
}
