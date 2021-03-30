// +build goleveldb

package goleveldb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func TestGoLevelDBNewDB(t *testing.T) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)
}

func TestGoLevelDBNewDBAndNewReadonlyDB(t *testing.T) {
	// Test we can't open the db twice for writing
	name, dir := dbtest.NewTestName("goleveldb")
	wr1, err := NewDB(name, dir)
	defer dbtest.CleanupDB(wr1, name, dir)
	require.NoError(t, err)
	_, err = NewDB(name, dir)
	require.Error(t, err)
	wr1.Close() // Close the db to release the lock

	// Test we can open the db twice for reading only
	ro1, err := NewDBWithOpts(name, dir, &opt.Options{ReadOnly: true})
	require.NoError(t, err)
	defer ro1.Close()
	ro2, err := NewDBWithOpts(name, dir, &opt.Options{ReadOnly: true})
	require.NoError(t, err)
	defer ro2.Close()
}

func TestGoLevelDBStats(t *testing.T) {
	name, dir := dbtest.NewTestName("cleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	assert.NotEmpty(t, db.Stats())
}

func TestGoLevelDBIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBIterator(t, db)
}

func TestGoLevelDBEmptyIterator(t *testing.T) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBEmptyIterator(t, db)
}

func TestGoLevelDBBatch(t *testing.T) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(t, err)

	dbtest.TestDBBatch(t, db)
}

func BenchmarkGoLevelDBRangeScans1M(b *testing.B) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkGoLevelDBRangeScans10M(b *testing.B) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkGoLevelDBRandomReadsWrites(b *testing.B) {
	name, dir := dbtest.NewTestName("goleveldb")
	db, err := NewDB(name, dir)
	defer dbtest.CleanupDB(db, name, dir)
	require.NoError(b, err)

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
