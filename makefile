GOTOOLS = github.com/golangci/golangci-lint/cmd/golangci-lint
PACKAGES=$(shell go list ./...)
INCLUDE = -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf

export GO111MODULE = on

all: lint test

### go tests
## By default this will only test memdb & goleveldb
test:
	@echo "--> Running go test"
	@go test $(PACKAGES) -tags memdb,goleveldb -v

test-memdb:
	@echo "--> Running go test"
	@go test ./memdb/... -tags memdb -v

test-goleveldb:
	@echo "--> Running go test"
	@go test ./goleveldb/... -tags goleveldb -v

test-cleveldb:
	@echo "--> Running go test"
	@go test ./cleveldb/... -tags cleveldb -v

test-rocksdb:
	@echo "--> Running go test"
	@go test ./rocksdb/... -tags rocksdb -v

test-boltdb:
	@echo "--> Running go test"
	@go test ./boltdb/... -tags boltdb -v

test-badgerdb:
	@echo "--> Running go test"
	@go test ./badgerdb/... -tags badgerdb -v

test-prefixdb:
	@echo "--> Running go test"
	@go test ./prefixdb/... -tags prefixdb -v

test-remotedb:
	@echo "--> Running go test"
	@go test ./remotedb/... -tags goleveldb,remotedb -v

test-all:
	@echo "--> Running go test"
	@go test $(PACKAGES) -tags memdb,goleveldb,cleveldb,boltdb,rocksdb,badgerdb,prefixdb,remotedb -v

test-all-docker:
	@echo "--> Running go test"
	@docker run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-tm-db-testing go test $(PACKAGES) -tags memdb,goleveldb,cleveldb,boltdb,rocksdb,badgerdb,prefixdb,remotedb -v
.PHONY: test-all-docker

bench:
	@go test -bench=. $(PACKAGES) -tags memdb,goleveldb

bench-memdb:
	@go test -bench=. ./memdb/... -tags memdb

bench-goleveldb:
	@go test -bench=. ./goleveldb/... -tags goleveldb

bench-cleveldb:
	@go test -bench=. ./cleveldb/... -tags cleveldb

bench-rocksdb:
	@go test -bench=. ./rocksdb/... -tags rocksdb

bench-boltdb:
	@go test -bench=. ./boltdb/... -tags boltdb

bench-badgerdb:
	@go test -bench=. ./badgerdb/... -tags badgerdb

bench-prefixdb:
	@go test -bench=. ./prefixdb/... -tags prefixdb

bench-remotedb:
	@go test -bench=. ./remotedb/... -tags goleveldb,remotedb

bench-all:
	@go test -bench=. $(PACKAGES) -tags memdb,goleveldb,cleveldb,boltdb,rocksdb,badgerdb,prefixdb,remotedb

bench-all-docker:
	@docker run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-tm-db-testing go test -bench=. $(PACKAGES) -tags memdb,goleveldb,cleveldb,boltdb,rocksdb,badgerdb,prefixdb,remotedb
.PHONY: bench-all-docker

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
.PHONY: lint

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*.pb.go' -not -name '*pb_test.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*"  -not -name '*.pb.go' -not -name '*pb_test.go' | xargs goimports -w
.PHONY: format

tools:
	go get -v $(GOTOOLS)

# generates certificates for TLS testing in remotedb
gen_certs: clean_certs
	certstrap init --common-name "tendermint.com" --passphrase ""
	certstrap request-cert --common-name "remotedb" -ip "127.0.0.1" --passphrase ""
	certstrap sign "remotedb" --CA "tendermint.com" --passphrase ""
	mv out/remotedb.crt remotedb/test.crt
	mv out/remotedb.key remotedb/test.key
	rm -rf out

clean_certs:
	rm -f db/remotedb/test.crt
	rm -f db/remotedb/test.key

%.pb.go: %.proto
	## If you get the following error,
	## "error while loading shared libraries: libprotobuf.so.14: cannot open shared object file: No such file or directory"
	## See https://stackoverflow.com/a/25518702
	## Note the $< here is substituted for the %.proto
	## Note the $@ here is substituted for the %.pb.go
	protoc $(INCLUDE) $< --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,plugins=grpc:.


protoc_remotedb: remotedb/proto/defs.pb.go	
