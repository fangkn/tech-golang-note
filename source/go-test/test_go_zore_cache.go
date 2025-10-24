package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
)

func main() {

	cache, err := collection.NewCache(time.Second*1, collection.WithName("any"))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)

	cache.SetWithExpire("first", "first element1111", time.Second*5)
	time.Sleep(time.Second * 3)

	v, exist := cache.Get("first")
	if exist {
		fmt.Println("first element:", v)
	} else {
		fmt.Println("first element expired")
	}
}
