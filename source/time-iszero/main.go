package main

import (
	"fmt"
	"time"
)

func main() {

	applyTime := time.Unix(0, 0)
	fmt.Println(applyTime)
	fmt.Println(applyTime.IsZero())
	fmt.Println(applyTime == time.Time{})
	fmt.Println(applyTime == time.Unix(0, 0))

	applyTime2 := time.Time{}
	fmt.Println(applyTime2.IsZero())

}
