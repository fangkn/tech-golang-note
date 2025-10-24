package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var i int32 = 100

func add() {

	atomic.AddInt32(&i, 1) // cas compare and swap old new
	fmt.Printf("i++: %v\n", i)
}

func sub() {

	atomic.AddInt32(&i, -1)

	fmt.Printf("i-- %v\n", i)
}

func main() {

	for i := 0; i < 100; i++ {

		go add()

		go sub()
	}

	time.Sleep(time.Second * 3)
	fmt.Printf("end i: %v\n", i)
}
