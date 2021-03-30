// +build badgerdb

package badgerdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/require"
)

func TestBadgerDBBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

// TODO implement badgerDB.Status()
// func TestBadgerDBStats(t *testing.T) {
// 	name, dir := dbtest.NewTestName("badgerdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(t, err)
//
// 	assert.NotEmpty(t, db.Stats())
// }

func TestBadgerDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestBadgerDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestBadgerDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

func BenchmarkBadgerDBRangeScans1M(b *testing.B) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkBadgerDBRangeScans10M(b *testing.B) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkBadgerDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
