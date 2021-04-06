// +build prefixdb

package prefixdb_test

import (
	"testing"

	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/line/tm-db/v2/memdb"
	"github.com/line/tm-db/v2/prefixdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockDBWithStuff(t *testing.T) tmdb.DB {
	db := memdb.NewDB()
	// Under "key" prefix
	require.NoError(t, db.Set([]byte("key"), []byte("value")))
	require.NoError(t, db.Set([]byte("key1"), []byte("value1")))
	require.NoError(t, db.Set([]byte("key2"), []byte("value2")))
	require.NoError(t, db.Set([]byte("key3"), []byte("value3")))
	require.NoError(t, db.Set([]byte("something"), []byte("else")))
	require.NoError(t, db.Set([]byte("k"), []byte("val")))
	require.NoError(t, db.Set([]byte("ke"), []byte("valu")))
	require.NoError(t, db.Set([]byte("kee"), []byte("valuu")))
	return db
}

func TestPrefixDBSimple(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	dbtest.Value(t, pdb, []byte("key"), nil)
	dbtest.Value(t, pdb, []byte("key1"), nil)
	dbtest.Value(t, pdb, []byte("1"), []byte("value1"))
	dbtest.Value(t, pdb, []byte("key2"), nil)
	dbtest.Value(t, pdb, []byte("2"), []byte("value2"))
	dbtest.Value(t, pdb, []byte("key3"), nil)
	dbtest.Value(t, pdb, []byte("3"), []byte("value3"))
	dbtest.Value(t, pdb, []byte("something"), nil)
	dbtest.Value(t, pdb, []byte("k"), nil)
	dbtest.Value(t, pdb, []byte("ke"), nil)
	dbtest.Value(t, pdb, []byte("kee"), nil)
}

func TestPrefixDBIterator1(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	itr, err := pdb.Iterator(nil, nil)
	require.NoError(t, err)
	dbtest.Item(t, itr, []byte("1"), []byte("value1"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("2"), []byte("value2"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("3"), []byte("value3"))
	dbtest.Next(t, itr, false)
	dbtest.Invalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator1(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	itr, err := pdb.ReverseIterator(nil, nil)
	require.NoError(t, err)
	dbtest.Item(t, itr, []byte("3"), []byte("value3"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("2"), []byte("value2"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("1"), []byte("value1"))
	dbtest.Next(t, itr, false)
	dbtest.Invalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator5(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	itr, err := pdb.ReverseIterator([]byte("1"), nil)
	require.NoError(t, err)
	dbtest.Item(t, itr, []byte("3"), []byte("value3"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("2"), []byte("value2"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("1"), []byte("value1"))
	dbtest.Next(t, itr, false)
	dbtest.Invalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator6(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	itr, err := pdb.ReverseIterator([]byte("2"), nil)
	require.NoError(t, err)
	dbtest.Item(t, itr, []byte("3"), []byte("value3"))
	dbtest.Next(t, itr, true)
	dbtest.Item(t, itr, []byte("2"), []byte("value2"))
	dbtest.Next(t, itr, false)
	dbtest.Invalid(t, itr)
	itr.Close()
}

func TestPrefixDBReverseIterator7(t *testing.T) {
	db := mockDBWithStuff(t)
	pdb := prefixdb.NewDB(db, []byte("key"))

	itr, err := pdb.ReverseIterator(nil, []byte("2"))
	require.NoError(t, err)
	dbtest.Item(t, itr, []byte("1"), []byte("value1"))
	dbtest.Next(t, itr, false)
	dbtest.Invalid(t, itr)
	itr.Close()
}

func TestPrefixDBNewDB(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	require.Panics(t, func() {
		prefixdb.NewDB(db, []byte{})
	})

	require.Panics(t, func() {
		prefixdb.NewDB(db, nil)
	})
}

func TestPrefixDBStats(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	assert.NotEmpty(t, pdb.Stats())
}

func TestPrefixDBIterator(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestDBIterator(t, pdb)
}

func TestPrefixDBIteratorNoWrites(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestDBIteratorNoWrites(t, pdb)
}

func TestPrefixDBEmptyIterator(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestDBEmptyIterator(t, pdb)
}

func TestPrefixDBPrefixIterator(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestDBPrefixIterator(t, pdb)
}

func TestPrefixDBPrefixIteratorNoMatchNil(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestPrefixIteratorNoMatchNil(t, pdb)
}

func TestPrefixDBPrefixIteratorNoMatch1(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestPrefixIteratorNoMatch1(t, pdb)
}

func TestPrefixDBPrefixIteratorNoMatch2(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestPrefixIteratorNoMatch2(t, pdb)
}

func TestPrefixDBPrefixIteratorMatch1(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestPrefixIteratorMatch1(t, pdb)
}

func TestPrefixDBPrefixIteratorMatches1N(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestPrefixIteratorMatches1N(t, pdb)
}

func TestPrefixDBBatch(t *testing.T) {
	db := memdb.NewDB()
	require.NotNil(t, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(t, pdb)
	defer pdb.Close()

	dbtest.TestDBBatch(t, pdb)
}

func BenchmarkPrefixDBRangeScans1M(b *testing.B) {
	db := memdb.NewDB()
	require.NotNil(b, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(b, pdb)
	defer pdb.Close()

	dbtest.BenchmarkRangeScans(b, pdb, int64(1e6))
}

func BenchmarkPrefixDBRangeScans10M(b *testing.B) {
	db := memdb.NewDB()
	require.NotNil(b, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(b, pdb)
	defer pdb.Close()

	dbtest.BenchmarkRangeScans(b, pdb, int64(10e6))
}

func BenchmarkPrefixDBRandomReadsWrites(b *testing.B) {
	db := memdb.NewDB()
	require.NotNil(b, db)
	defer db.Close()
	pdb := prefixdb.NewDB(db, []byte("key"))
	require.NotNil(b, pdb)
	defer pdb.Close()

	dbtest.BenchmarkRandomReadsWrites(b, pdb)
}
