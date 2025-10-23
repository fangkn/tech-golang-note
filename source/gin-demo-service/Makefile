


all: release

.PHONY: debug
debug:
	go vet ./...
	go build -ldflags 

.PHONY: release
release:
	go vet ./...
	GOWORK=off go build -ldflags '-w -s'

.PHONY: docs
docs:
	swag init
	swagger generate markdown -f docs/swagger.yaml --output Readme.md
