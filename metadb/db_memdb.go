package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/memdb"
)

func memdbDBCreator(name, dir string) (tmdb.DB, error) {
	return memdb.NewDB(), nil
}

func init() { registerDBCreator(MemDBBackend, memdbDBCreator, false) }
