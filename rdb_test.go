//go:build rocksdb
// +build rocksdb

// Avoiding duplicate codes by lint, this implementation re-ordered functions
package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(10e6))
}

func TestRDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, RDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	_, ok := db.(*RDB)
	assert.True(t, ok)
}

func TestRDBStats(t *testing.T) {
	db, dir, name := newDB(t, RDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	assert.NotEmpty(t, db.Stats())
}

func BenchmarkRDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRandomReadsWrites(b, db)
}
