BUNDLE_PATH := "tmp/bundle"

all: test build

build: 
	go build -o bin/clickhouse-evm ./cmd

bundle: build 
	mkdir -p ${BUNDLE_PATH}
	mkdir -p ${BUNDLE_PATH}/etc/clickhouse-server
	mkdir -p ${BUNDLE_PATH}/var/lib/clickhouse/user_defined
	mkdir -p ${BUNDLE_PATH}/var/lib/clickhouse/user_scripts
	cp bin/clickhouse-evm ${BUNDLE_PATH}/var/lib/clickhouse/user_scripts/
	cp config/*_function.*ml ${BUNDLE_PATH}/etc/clickhouse-server/
	cp sql/function_*.sql ${BUNDLE_PATH}/var/lib/clickhouse/user_defined/
	COPYFILE_DISABLE=1 tar --no-xattr -cvzf ${BUNDLE_PATH}/../bundle.tar.gz -C ${BUNDLE_PATH} .

test:
	go test -v $(shell go list ./... | grep -v /e2e)

e2e-test:
	go test -v ./e2e

clean:
	rm -rf bin
	rm -rf ${BUNDLE_PATH}/../bundle.tar.gz ${BUNDLE_PATH} ${BUNDLE_PATH}/../bundle.tar.gz
