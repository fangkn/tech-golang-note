package main

import (
	"fmt"
	"os"

	"github.com/iafan/cwalk"
)

func walkFunc(path string, info os.FileInfo, err error) error {
	fmt.Printf("path: %s\n", path)

	return nil
}

func main() {
	err := cwalk.Walk("c:\\kane-fang\\doc\\", walkFunc)
	if err != nil {
		fmt.Println(err)
	}
}
