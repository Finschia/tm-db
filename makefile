GOTOOLS = github.com/golangci/golangci-lint/cmd/golangci-lint
PACKAGES=$(shell go list ./...)
INCLUDE = -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf

# Setup
# See: LevelDB: https://github.com/jmhodges/levigo/blob/master/README.md
# See: RocksDB: https://github.com/line/gorocksdb/blob/main/README.md
CLEVELDB_DIR=$(shell pwd)/leveldb
ROCKSDB_DIR=$(shell pwd)/rocksdb
CGO_CFLAGS=-I$(CLEVELDB_DIR)/include -I$(ROCKSDB_DIR)/include
CGO_LDFLAGS=-L$(CLEVELDB_DIR) -L$(ROCKSDB_DIR) -lleveldb -lrocksdb -lm -lstdc++ $(shell awk '/PLATFORM_LDFLAGS/ {sub("PLATFORM_LDFLAGS=", ""); print}' < $(ROCKSDB_DIR)/make_config.mk)

DOCKER_NAME=tm-db-testing
DOCKER_IMAGE=line/$(DOCKER_NAME)

export GO111MODULE = on

all: lint test

### go tests
## By default this will only test memdb & goleveldb
test:
	@echo "--> Running go test"
	@go test $(PACKAGES) -v

test-cleveldb: build-cleveldb
	@echo "--> Running go test"
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test $(PACKAGES) -tags cleveldb -v

test-rocksdb: build-rocksdb
	@echo "--> Running go test"
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test $(PACKAGES) -tags rocksdb -v

test-boltdb:
	@echo "--> Running go test"
	@go test $(PACKAGES) -tags boltdb -v

test-badgerdb:
	@echo "--> Running go test"
	@go test $(PACKAGES) -tags badgerdb -v

test-all: build-cleveldb build-rocksdb
	@echo "--> Running go test"
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test $(PACKAGES) -tags cleveldb,rocksdb,boltdb,badgerdb -v

test-all-docker:
	@echo "--> Running go test"
	@docker run --rm -e CGO_LDFLAGS="-lrocksdb" -v $(CURDIR):/workspace --workdir /workspace $(DOCKER_IMAGE) \
	go test $(PACKAGES) -tags cleveldb,rocksdb,boltdb,badgerdb -v
.PHONY: test-all-docker

bench:
	@go test -bench=. $(PACKAGES)

bench-cleveldb: build-cleveldb
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test -bench=. $(PACKAGES) -tags cleveldb

bench-rocksdb: build-rocksdb
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test -bench=. $(PACKAGES) -tags rocksdb

bench-boltdb:
	@go test -bench=. $(PACKAGES) -tags boltdb

bench-badgerdb:
	@go test -bench=. $(PACKAGES) -tags badgerdb

bench-all: build-cleveldb build-rocksdb
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go test -bench=. $(PACKAGES) -tags cleveldb,rocksdb,boltdb,badgerdb

bench-all-docker:
	@docker run --rm -e CGO_LDFLAGS="-lrocksdb" -v $(CURDIR):/workspace --workdir /workspace $(DOCKER_IMAGE) \
	go test -bench=. $(PACKAGES) -tags cleveldb,rocksdb,boltdb,badgerdb
.PHONY: bench-all-docker

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
.PHONY: lint

lint-all: build-cleveldb build-rocksdb
	@echo "--> Running linter"
	@CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	golangci-lint run --build-tags "cleveldb,rocksdb,boltdb,badgerdb"
	@go mod verify
.PHONY: lint-all

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*.pb.go' -not -name '*pb_test.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*"  -not -name '*.pb.go' -not -name '*pb_test.go' | xargs goimports -w
.PHONY: format

tools:
	go get -v $(GOTOOLS)

build-cleveldb:
	@if [ ! -e $(CLEVELDB_DIR) ]; then \
		sh ./contrib/get_cleveldb.sh; \
	fi
	@if [ ! -e $(CLEVELDB_DIR)/libcleveldb.a ]; then \
		cd $(CLEVELDB_DIR) && make; \
	fi
.PHONY: build-cleveldb

build-rocksdb:
	@if [ ! -e $(ROCKSDB_DIR) ]; then \
		sh ./contrib/get_rocksdb.sh; \
	fi
	@if [ ! -e $(ROCKSDB_DIR)/librocksdb.a ]; then \
		cd $(ROCKSDB_DIR) && make -j4 static_lib; \
	fi
.PHONY: build-rocksdb

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

build-local-docker:
	## If you met the compile error, it might be the shortage of memory on docker with your machine
	## Please check your Docker Desktop and its settings of resources
	## for mac: https://docs.docker.com/desktop/mac/#resources
	## for windows: https://docs.docker.com/desktop/windows/#resources
	@docker build --rm --progress plain --tag="$(DOCKER_IMAGE)" -f ./tools/Dockerfile .
.PHONY: build-local-docker

bash-local-docker:
	## If you want to `golangci-lint` in Docker,
	## should install `golangci-lint` after Docker container start with bash and
	## should be the same version of lint.yml.
	## See: .github/lint.yml
	##
	## example: golangci-lint
	##
	## curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1
	## golangci-lint run --build-tags "cleveldb,rocksdb,boltdb,badgerdb"
	@docker run --rm -it -e CGO_LDFLAGS="-lrocksdb" -v $(CURDIR):/workspace --workdir /workspace $(DOCKER_IMAGE) bash
.PHONY: bash-local-docker
