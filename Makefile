
.PHONY: all
all: build test

.PHONY: build
build:
	go run ./script/build.go build

.PHONY: test
test:
	go test ./...
	go vet ./...
