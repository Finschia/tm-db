// +build badgerdb

package badgerdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/require"
)

func TestBadgerDBNewDB(t *testing.T) {
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

func TestBadgerDBIteratorNoWrites(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIteratorNoWrites(t, db)
}

func TestBadgerDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestBadgerDBPrefixIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBPrefixIterator(t, db)
}

func TestBadgerDBPrefixIteratorNoMatchNil(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatchNil(t, db)
}

func TestBadgerDBPrefixIteratorNoMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch1(t, db)
}

func TestBadgerDBPrefixIteratorNoMatch2(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch2(t, db)
}

func TestBadgerDBPrefixIteratorMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatch1(t, db)
}

func TestBadgerDBPrefixIteratorMatches1N(t *testing.T) {
	name, dir := dbtest.NewTestName("badgerdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatches1N(t, db)
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
