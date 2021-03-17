PKG_LIST := $(shell go list ./...)

all: build

clean:
	go clean
	rm -rf bin

test: clean
	test -z '$(shell gofmt -l .)'
	go vet ./...
	go test ./... -v

build: test
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/bambus cmd/*
