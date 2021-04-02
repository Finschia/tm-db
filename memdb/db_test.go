// +build memdb

package memdb

import (
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemDBNewDB(t *testing.T) {
	db := NewDB()
	require.NotNil(t, db)
	db.Close()
}

func TestMemDBStats(t *testing.T) {
	db := NewDB()
	defer db.Close()

	assert.NotEmpty(t, db.Stats())
}

func TestMemDBIterator(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestDBIterator(t, db)
}

func TestMemDBEmptyIterator(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestDBEmptyIterator(t, db)
}

func TestMemDBPrefixIterator(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestDBPrefixIterator(t, db)
}

func TestMemDBPrefixIteratorNoMatchNil(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestPrefixIteratorNoMatchNil(t, db)
}

func TestMemDBPrefixIteratorNoMatch1(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestPrefixIteratorNoMatch1(t, db)
}

func TestMemDBPrefixIteratorNoMatch2(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestPrefixIteratorNoMatch2(t, db)
}

func TestMemDBPrefixIteratorMatch1(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestPrefixIteratorMatch1(t, db)
}

func TestMemDBPrefixIteratorMatches1N(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestPrefixIteratorMatches1N(t, db)
}

func TestMemDBBatch(t *testing.T) {
	db := NewDB()
	defer db.Close()

	dbtest.TestDBBatch(t, db)
}

func BenchmarkMemDBRangeScans1M(b *testing.B) {
	db := NewDB()
	defer db.Close()

	dbtest.BenchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkMemDBRangeScans10M(b *testing.B) {
	db := NewDB()
	defer db.Close()

	dbtest.BenchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkMemDBRandomReadsWrites(b *testing.B) {
	db := NewDB()
	defer db.Close()

	dbtest.BenchmarkRandomReadsWrites(b, db)
}
