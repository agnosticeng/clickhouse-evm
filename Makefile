all: test build

build: 
	go build -o bin/agnostic-clickhouse-udf ./cmd

test:
	go test -v ./...

clean:
	rm -rf bin
