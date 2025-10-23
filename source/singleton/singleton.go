package singleton

import (
	"fmt"
	"sync"
)

type Singleton struct {
	Value int
}

var (
	instance *Singleton
	once     sync.Once
)

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating Singleton instance")
		instance = &Singleton{Value: 100}
	})
	return instance
}
