// +build cleveldb

package cleveldb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLevelDBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

func TestCLevelDBStats(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func TestCLevelDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestCLevelDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

func BenchmarkCLevelDBRangeScans1M(b *testing.B) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkCLevelDBRangeScans10M(b *testing.B) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkCLevelDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
