// +build rocksdb

package rocksdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRocksDBBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

func TestRocksDBStats(t *testing.T) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func TestRocksDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestRocksDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

func BenchmarkRocksDBRangeScans1M(b *testing.B) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRocksDBRangeScans10M(b *testing.B) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkRocksDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("rocksdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
