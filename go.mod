module github.com/line/tm-db/v2

go 1.15

require (
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/gogo/protobuf v1.3.2
	github.com/google/btree v1.0.0
	github.com/jmhodges/levigo v1.0.0
	github.com/line/gorocksdb v0.0.0-20210405045146-6e1a987c6552
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	go.etcd.io/bbolt v1.3.5
	google.golang.org/grpc v1.35.0
)

replace github.com/line/gorocksdb => ../gorocksdb
