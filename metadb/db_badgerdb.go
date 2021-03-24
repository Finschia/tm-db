// +build badgerdb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/badgerdb"
)

func badgerDBCreator(name, dir string) (tmdb.DB, error) {
	return badgerdb.NewDB(name, dir)
}

func init() { registerDBCreator(BadgerDBBackend, badgerDBCreator, true) }
