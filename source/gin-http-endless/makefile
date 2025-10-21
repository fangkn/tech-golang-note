BINARY := gin-http-endless

.PHONY: build run clean tidy

build:
	@echo "Building $(BINARY)..."
	CGO_ENABLED=0 go build -o $(BINARY) .

run:
	@echo "Running..."
	go run .

tidy:
	@echo "Tidying modules..."
	go mod tidy

clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
