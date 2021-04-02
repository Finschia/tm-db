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

func TestCLevelDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestCLevelDBPrefixIteratorNoMatchNil(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatchNil(t, db)
}

func TestCLevelDBPrefixIteratorNoMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch1(t, db)
}

func TestCLevelDBPrefixIteratorNoMatch2(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorNoMatch2(t, db)
}

func TestCLevelDBPrefixIteratorMatch1(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatch1(t, db)
}

func TestCLevelDBPrefixIteratorMatches1N(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestPrefixIteratorMatches1N(t, db)
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
