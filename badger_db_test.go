//go:build badgerdb
// +build badgerdb

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerDBStats(t *testing.T) {
	db, dir, name := newDB(t, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	assert.Nil(t, db.Stats()) // Not implement
}

func TestBadgerDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	_, ok := db.(*BadgerDB)
	assert.True(t, ok)
}

func BenchmarkBadgerDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRandomReadsWrites(b, db)
}

func BenchmarkBadgerDBParallelRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkParallelRandomReadsWrites(b, db)
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBadgerDBRandomReadsWrites
func TempBenchmarkBadgerDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(1e6))
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBadgerDBRandomReadsWrites
func TempBenchmarkBadgerDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, BadgerDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(10e6))
}
