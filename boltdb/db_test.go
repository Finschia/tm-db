// +build boltdb

package boltdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoltDBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

func TestBoltDBStats(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func TestBoltDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestBoltDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestBoltDBPrefixIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBPrefixIterator(t, db)
}

func TestBoltDBPrefixIteratorNoMatchNil(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatchNil(t, db)
}

// TODO bolt does not support concurrent writes while iterating
// func TestBoltDBPrefixIteratorNoMatch1(t *testing.T) {
// 	name, dir := dbtest.NewTestName("boltdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(t, err)
//
// 	dbtest.TestPrefixIteratorNoMatch1(t, db)
// }

func TestBoltDBPrefixIteratorNoMatch2(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch2(t, db)
}

func TestBoltDBPrefixIteratorMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatch1(t, db)
}

func TestBoltDBPrefixIteratorMatches1N(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatches1N(t, db)
}

func TestBoltDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

// TODO fix stall
// func BenchmarkBoltDBRangeScans1M(b *testing.B) {
// 	name, dir := dbtest.NewTestName("boltdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(b, err)
//
// 	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
// }

// TODO fix stall
// func BenchmarkBoltDBRangeScans10M(b *testing.B) {
// 	name, dir := dbtest.NewTestName("boltdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(b, err)
//
// 	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
// }

func BenchmarkBoltDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("boltdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
