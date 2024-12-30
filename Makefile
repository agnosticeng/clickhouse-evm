BUNDLE_PATH := "tmp/bundle"

all: test build

build: 
	go build -o bin/agnostic-clickhouse-udf ./cmd

bundle: build 
	mkdir -p ${BUNDLE_PATH}
	mkdir -p ${BUNDLE_PATH}/etc/clickhouse-server
	mkdir -p ${BUNDLE_PATH}/var/lib/clickhouse/user_defined
	mkdir -p ${BUNDLE_PATH}/var/lib/clickhouse/user_scripts
	cp bin/agnostic-clickhouse-udf ${BUNDLE_PATH}/var/lib/clickhouse/user_scripts/
	cp config/*_function.*ml ${BUNDLE_PATH}/etc/clickhouse-server/
	cp sql/function_*.sql ${BUNDLE_PATH}/var/lib/clickhouse/user_defined/
	tar -cvzf ${BUNDLE_PATH}/../bundle.tar.gz -C ${BUNDLE_PATH} .

test:
	go test -v ./...

clean:
	rm -rf bin
	rm -rf ${BUNDLE_PATH} 
