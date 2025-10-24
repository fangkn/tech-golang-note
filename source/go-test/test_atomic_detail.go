package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func add() {

	var i int32 = 100
	atomic.AddInt32(&i, 1)
	fmt.Printf("i++: %v\n", i)
}

func main() {

	for i := 0; i < 100; i++ {

		go add()
	}

	time.Sleep(time.Second * 3)
	//fmt.Printf("end i: %v\n", i)
}
