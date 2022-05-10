//go:build badgerdb
// +build badgerdb

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadgerDBNewDB(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BadgerDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(t, err)

	_, ok := db.(*BadgerDB)
	assert.True(t, ok)
}

func TestBadgerDBStats(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BadgerDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(t, err)

	assert.Nil(t, db.Stats()) // Not implement
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBadgerDBRandomReadsWrites
func TempBenchmarkBadgerDBRangeScans1M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BadgerDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(1e6))
}

// Cannot work well since the data setup time is long (10min over)
// See the read/write performance: BenchmarkBadgerDBRandomReadsWrites
func TempBenchmarkBadgerDBRangeScans10M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, BadgerDBBackend, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkBadgerDBRandomReadsWrites(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewBadgerDB(name, dir)
	defer closeDBWithCleanupDBDir(db, dir, name)
	require.NoError(b, err)

	benchmarkRandomReadsWrites(b, db)
}
