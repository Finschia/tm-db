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

func TestRDBNewDB(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RDBBackend, dir)
	defer cleanupDBDir(dir, name)
	require.NoError(t, err)

	_, ok := db.(*RDB)
	assert.True(t, ok)
}

func TestRDBStats(t *testing.T) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RDBBackend, dir)
	defer cleanupDBDir(dir, name)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func BenchmarkRDBRangeScans1M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RDBBackend, dir)
	defer cleanupDBDir(dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRDBRangeScans10M(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RDBBackend, dir)
	defer cleanupDBDir(dir, name)
	require.NoError(b, err)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkRDBRandomReadsWrites(b *testing.B) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, RDBBackend, dir)
	defer cleanupDBDir(dir, name)
	require.NoError(b, err)

	benchmarkRandomReadsWrites(b, db)
}
