//go:build boltdb
// +build boltdb

// Avoiding duplicate codes by lint, this implementation re-ordered functions
package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, BoltDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	_, ok := db.(*BoltDB)
	assert.True(t, ok)
}

func TestBoltDBStats(t *testing.T) {
	db, dir, name := newDB(t, BoltDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	assert.NotEmpty(t, db.Stats())
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBoltDBRandomReadsWrites
func TempBenchmarkBoltDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, BoltDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(1e6))
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBoltDBRandomReadsWrites
func TempBenchmarkBoltDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, BoltDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkBoltDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, BoltDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRandomReadsWrites(b, db)
}
