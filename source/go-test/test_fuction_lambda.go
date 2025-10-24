package main

import "fmt"

func main() {

	// max := func(a int, b int) int {
	// 	if a > b {
	// 		return a
	// 	} else {
	// 		return b
	// 	}
	// }
	//
	// r := max(1, 2)
	// fmt.Printf("r: %v\n", r)

	r := func(a int, b int) int {
		if a > b {
			return a
		} else {
			return b
		}
	}(1, 2)

	fmt.Printf("r: %v\n", r)

}
