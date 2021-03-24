// +build rocksdb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/rocksdb"
)

func rocksDBCreator(name, dir string) (tmdb.DB, error) {
	return rocksdb.NewDB(name, dir)
}

func init() { registerDBCreator(RocksDBBackend, rocksDBCreator, true) }
