//go:build rocksdb
// +build rocksdb

package rdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRDBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

func TestRDBStats(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func TestRDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestRDBIteratorNoWrites(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIteratorNoWrites(t, db)
}

func TestRDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestRDBPrefixIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBPrefixIterator(t, db)
}

func TestRDBPrefixIteratorNoMatchNil(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatchNil(t, db)
}

func TestRDBPrefixIteratorNoMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch1(t, db)
}

func TestRDBPrefixIteratorNoMatch2(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch2(t, db)
}

func TestRDBPrefixIteratorMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatch1(t, db)
}

func TestRDBPrefixIteratorMatches1N(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatches1N(t, db)
}

func TestRDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

func BenchmarkRDBRangeScans1M(b *testing.B) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkRDBRangeScans10M(b *testing.B) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkRDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("rdb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
