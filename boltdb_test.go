//go:build boltdb
// +build boltdb

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoltDBNewDB(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BoltDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(t, err)

	_, ok := db.(*BoltDB)
	assert.True(t, ok)
}

func TestBoltDBStats(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BoltDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBoltDBRandomReadsWrites
func TempBenchmarkBoltDBRangeScans1M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BoltDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(1e6))
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBoltDBRandomReadsWrites
func TempBenchmarkBoltDBRangeScans10M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BoltDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkBoltDBRandomReadsWrites(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewBoltDB(name, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRandomReadsWrites(b, db)
}
