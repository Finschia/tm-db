package metadb

import (
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

// IteratePrefix is a convenience function for iterating over a key domain
// restricted by prefix.
func IteratePrefix(db tmdb.DB, prefix []byte) (tmdb.Iterator, error) {
	var start, end []byte
	if len(prefix) == 0 {
		start = nil
		end = nil
	} else {
		start = util.Cp(prefix)
		end = util.CpIncr(prefix)
	}
	itr, err := db.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return itr, nil
}
