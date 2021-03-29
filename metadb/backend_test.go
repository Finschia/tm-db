package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/memdb"
	"github.com/line/tm-db/v2/prefixdb"
)

// TODO
// Register a test backend for PrefixDB as well, with some unrelated junk data
func init() {
	// nolint: errcheck
	registerDBCreator("prefixdb", func(name, dir string) (tmdb.DB, error) {
		mdb := memdb.NewDB()
		mdb.Set([]byte("a"), []byte{1})
		mdb.Set([]byte("b"), []byte{2})
		mdb.Set([]byte("t"), []byte{20})
		mdb.Set([]byte("test"), []byte{0})
		mdb.Set([]byte("u"), []byte{21})
		mdb.Set([]byte("z"), []byte{26})
		return prefixdb.NewDB(mdb, []byte("test/")), nil
	}, false)
}
