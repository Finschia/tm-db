// +build goleveldb

package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/goleveldb"
)

func golevelDBCreator(name, dir string) (tmdb.DB, error) {
	return goleveldb.NewDB(name, dir)
}

func init() { registerDBCreator(GoLevelDBBackend, golevelDBCreator, true) }
