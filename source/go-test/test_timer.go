package main

import (
	"fmt"
	"time"
)

func main() {

	/// -----------------
	// timer := time.NewTimer(time.Second * 2)
	// fmt.Printf("time.Now() : %v\n", time.Now())
	// t1 := <-timer.C // 阻塞的，指定的时间到了
	//
	// fmt.Printf("t1:%v\n", t1)

	fmt.Printf("time.Now() : %v\n", time.Now())
	timer := time.NewTimer(time.Second * 2)
	<-timer.C
	fmt.Printf("time.Now() : %v\n", time.Now())
}
