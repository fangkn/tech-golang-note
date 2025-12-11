.PHONY: run build clean test

# 运行程序
run:
	go run main.go

# 构建程序
build:
	go build -o validator-demo main.go

# 清理构建文件
clean:
	rm -f validator-demo

# 下载依赖
deps:
	go mod download
	go mod tidy

# 运行测试（如果有测试文件）
test:
	go test -v ./...

