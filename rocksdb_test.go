//go:build rocksdb
// +build rocksdb

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRocksDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, RDBBackend)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`

	_, ok := db.(*RocksDB)
	assert.True(t, ok)
}

func TestRocksDBStats(t *testing.T) {
	db, dir, name := newDB(t, RDBBackend)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`

	assert.NotEmpty(t, db.Stats())
}

func BenchmarkRocksDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRocksDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkRocksDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, RDBBackend)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`

	benchmarkRandomReadsWrites(b, db)
}
