package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGoLevelDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, GoLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkGoLevelDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, GoLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkGoLevelDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, GoLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRandomReadsWrites(b, db)
}

func TestGoLevelDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, GoLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	_, ok := db.(*GoLevelDB)
	assert.True(t, ok)
}

func TestGoLevelDBStats(t *testing.T) {
	db, dir, name := newDB(t, GoLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	assert.NotEmpty(t, db.Stats())
}
