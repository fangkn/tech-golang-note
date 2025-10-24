package main

import (
	"fmt"
	"os"
)

func main() {
	//s := os.Environ()
	//fmt.Printf("s: %v\n", s)

	// s = os.Getenv("GOPATH")
	s, b := os.LookupEnv("GOMODCACHE")
	if b {
		fmt.Printf("s: %v\n", s)
	}

}
