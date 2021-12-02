// +build rocksdb

package metadb

import (
	"fmt"
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/rdb"
)

func rdbCreator(name, dir string) (tmdb.DB, error) {
	fmt.Printf("XXX: rdb %s/%s\n", name, dir)
	return rdb.NewDB(name, dir)
}

func init() { registerDBCreator(RDBBackend, rdbCreator, true) }
