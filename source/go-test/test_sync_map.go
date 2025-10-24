package main

import (
	"fmt"
	"sync"
)

var scene sync.Map

func ListTasks() (tasks []string) {
	scene.Range(func(k, v interface{}) bool {
		tasks = append(tasks, k.(string))
		return true
	})
	return
}
func main() {

	// 将键值对保存到sync.Map
	scene.Store("greece", 97)
	scene.Store("london", 100)
	scene.Store("egypt", 200)

	// 从sync.Map中根据键取值
	fmt.Println(scene.Load("london"))

	// 根据键删除对应的键值对
	scene.Delete("london")

	// 遍历所有sync.Map中的键值对
	// scene.Range(func(k, v interface{}) bool {

	// 	fmt.Println("iterate:", k, v)
	// 	return true
	// })

	scenes := ListTasks()
	fmt.Println(scenes)

}
