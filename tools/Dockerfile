# This file defines the container image used to build and test tm-db in CI.
# The CI workflows use the latest tag of line/tm-db-testing
# built from these settings.
#
# The jobs defined in the Build & Push workflow will build and update the image
# when changes to this file are merged.  If you have other changes that require
# updates here, merge the changes here first and let the image get updated (or
# push a new version manually) before PRs that depend on them.

FROM golang:1.16-bullseye AS build

ENV LD_LIBRARY_PATH=/usr/local/lib

RUN apt-get update && apt-get install -y --no-install-recommends \
    libbz2-dev libgflags-dev libsnappy-dev libzstd-dev zlib1g-dev \
    make tar wget

FROM build AS install

COPY ./contrib/ ./contrib

# Install cleveldb
RUN \
  ./contrib/get_cleveldb.sh \
  && cd leveldb \
  && make \
  && cp -a out-static/lib* out-shared/lib* /usr/local/lib \
  && cd include \
  && cp -a leveldb /usr/local/include \
  && ldconfig

# Install Rocksdb
RUN \
  ./contrib/get_rocksdb.sh \
  && cd rocksdb \
  && DEBUG_LEVEL=0 make -j4 shared_lib \
  && make install-shared \
  && ldconfig

RUN rm -rf ./leveldb-*.tar.gz leveldb
RUN rm -rf ./rocksdb-*.tar.gz rocksdb
RUN rm -rf ./contrib

# Install golangci for CI
RUN go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

# Download dependency modules for CI
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
RUN rm -rf go.mod go.sum
