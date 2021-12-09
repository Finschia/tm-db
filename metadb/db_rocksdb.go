//go:build rocksdb
// +build rocksdb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/rdb"
	_ "github.com/line/tm-db/v2/rocksdb"
)

func rocksDBCreator(name, dir string) (tmdb.DB, error) {
	// TODO: use rdb instead of rocksdb for now as gorocksdb doesn't have low priority write option
	return rdb.NewDB(name, dir)
}

func init() { registerDBCreator(RocksDBBackend, rocksDBCreator, true) }
