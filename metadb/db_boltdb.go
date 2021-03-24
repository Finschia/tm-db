// +build boltdb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/boltdb"
)

func boltDBCreator(name, dir string) (tmdb.DB, error) {
	return boltdb.NewDB(name, dir)
}

func init() { registerDBCreator(BoltDBBackend, boltDBCreator, true) }
