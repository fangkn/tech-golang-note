package singleton

import (
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	const goroutineCount = 100
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	instances := make(chan *Singleton, goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			inst := GetInstance()
			instances <- inst
		}()
	}

	wg.Wait()
	close(instances)

	var first *Singleton
	for inst := range instances {
		if first == nil {
			first = inst
		} else if first != inst {
			t.Error("Different instances detected")
		}
	}
}
