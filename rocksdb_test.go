//go:build rocksdb
// +build rocksdb

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRocksDBNewDB(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RocksDBBackend, dir)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`
	require.NoError(t, err)

	_, ok := db.(*RocksDB)
	assert.True(t, ok)
}

func TestRocksDBStats(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RocksDBBackend, dir)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func BenchmarkRocksDBRangeScans1M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RocksDBBackend, dir)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRocksDBRangeScans10M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RocksDBBackend, dir)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkRocksDBRandomReadsWrites(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewRocksDB(name, dir)
	defer cleanupDBDir(dir, name) // Cannot use `closeDBWithCleanupDBDir`
	require.NoError(b, err)

	benchmarkRandomReadsWrites(b, db)
}
