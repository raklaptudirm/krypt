# build and test the code
.PHONY: all
all: build test

# run build script in script/build
.PHONY: build
build:
	go run ./script/build.go build

# test for inconsistencies in code
.PHONY: test
test:
	go test ./...
	go vet ./...

# get todo comments from code files
.PHONY: todo
todo:
	grep -r "TODO" .
