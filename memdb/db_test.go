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
