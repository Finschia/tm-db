// +build rocksdb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/rdb"
)

func rdbCreator(name, dir string) (tmdb.DB, error) {
	return rdb.NewDB(name, dir)
}

func init() { registerDBCreator(RDBBackend, rdbCreator, true) }
