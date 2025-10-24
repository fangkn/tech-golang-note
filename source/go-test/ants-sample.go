package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

var (
	taskWorkerPool *ants.PoolWithFunc
	wg             sync.WaitGroup
)

func doWork(w interface{}) {

	work := w.(string)
	fmt.Printf("work : %v\n", work)

}

func main() {

	runTimes := 10

	taskWorkerPool, _ = ants.NewPoolWithFunc(5, doWork, ants.WithNonblocking(true))

	defer taskWorkerPool.Release()

	for i := 0; i < runTimes; i++ {

		err := taskWorkerPool.Invoke(strconv.Itoa(i))

		fmt.Println(err)

	}

	time.Sleep(time.Second * 10)
	fmt.Printf("running goroutines: %d\n", taskWorkerPool.Running())

	return
}
