package main

import (
	"fmt"
	"sync"
)

var once sync.Once

func initialize() {
	fmt.Println("Initializing...")
}

func worker(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Worker %d attempting to initialize.\n", id)
	once.Do(initialize) // 只执行一次
	fmt.Printf("Worker %d completed.\n", id)
}

func main() {
	var wg sync.WaitGroup

	// 启动多个协程
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(&wg, i)
	}

	wg.Wait()
}
