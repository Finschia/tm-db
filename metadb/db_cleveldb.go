// +build cleveldb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/cleveldb"
)

func clevelDBCreator(name string, dir string) (tmdb.DB, error) {
	return cleveldb.NewDB(name, dir)
}

func init() { registerDBCreator(CLevelDBBackend, clevelDBCreator, false) }
