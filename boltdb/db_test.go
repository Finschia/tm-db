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

// TODO fix stall at closing db
// func TestBoltDBIterator(t *testing.T) {
// 	name, dir := dbtest.NewTestName("boltdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(t, err)
//
// 	dbtest.TestDBIterator(t, db)
// }

// TODO fix failure of test
// func TestBoltDBBatch(t *testing.T) {
// 	name, dir := dbtest.NewTestName("boltdb")
// 	db, err := NewDB(name, dir)
// 	defer dbtest.CleanupDB(db, name, dir)
// 	require.NoError(t, err)
//
// 	dbtest.TestDBBatch(t, db)
// }

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
