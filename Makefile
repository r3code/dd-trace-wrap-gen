export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

./bin:
	mkdir ./bin

./bin/golangci-lint: ./bin
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./...

test:
	 go test -race ./...

test-template:
	go run main.go -i TestInterface -o ./tests/interface_with_datadog_trace.go ./tests

all: build lint test


build:
	go build main.go

clean:
	rm -rf ./bin

