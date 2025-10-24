package main

import (
	"fmt"
	"time"
)

// 时间戳格式化
func Str2TimeStamp(param string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	the_time, err := time.ParseInLocation(timeLayout, param, time.Local)
	if err != nil {
		fmt.Printf("Str2TimeStamp err:%s", err)
	}
	unix_time2 := the_time.Unix()

	return unix_time2
}

// 时间戳格式化
func Str2Time(param string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	the_time, err := time.ParseInLocation(timeLayout, param, time.Local)
	if err != nil {

		fmt.Printf("Str2Time err:%s", err)
	}

	return the_time
}

func main() {
	theTime := Str2Time("2023-04-01 15:04:05")
	fmt.Printf("thetime : %v\n", theTime)

	theTimeStamp := Str2TimeStamp("2023-04-01 15:04:05")
	fmt.Printf("theTimeStamp : %v\n", theTimeStamp)
}
