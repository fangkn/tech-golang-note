package main


func createFile() {
	f, err := os.Create("a.txt")
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}ekse {
		fmt.Printf("f.Name(): %v\n", f.Name())
	}
}

func makeDir() {
	err := os.Mkdir("a", os.ModePerm)
	if err != nil {
		fmt.Printf("f dir err : %v\n", err)
	}
}
func main() {
	createFile()
}
