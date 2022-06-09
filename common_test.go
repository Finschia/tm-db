package db

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//----------------------------------------
// Helper functions.

func checkValue(t *testing.T, db DB, key []byte, valueWanted []byte) {
	valueGot, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, valueWanted, valueGot)
}

func checkValid(t *testing.T, itr Iterator, expected bool) {
	valid := itr.Valid()
	require.Equal(t, expected, valid)
}

func checkNext(t *testing.T, itr Iterator, expected bool) {
	itr.Next()
	// assert.NoError(t, err) TODO: look at fixing this
	valid := itr.Valid()
	require.Equal(t, expected, valid)
}

func checkNextPanics(t *testing.T, itr Iterator) {
	assert.Panics(t, func() { itr.Next() }, "checkNextPanics expected an error but didn't")
}

func checkItem(t *testing.T, itr Iterator, key []byte, value []byte) {
	v := itr.Value()

	k := itr.Key()

	assert.Exactly(t, key, k)
	assert.Exactly(t, value, v)
}

func checkInvalid(t *testing.T, itr Iterator) {
	checkValid(t, itr, false)
	checkKeyPanics(t, itr)
	checkValuePanics(t, itr)
	checkNextPanics(t, itr)
}

func checkKeyPanics(t *testing.T, itr Iterator) {
	assert.Panics(t, func() { itr.Key() }, "checkKeyPanics expected panic but didn't")
}

func checkValuePanics(t *testing.T, itr Iterator) {
	assert.Panics(t, func() { itr.Value() })
}

func newTempDB(t *testing.T, backend BackendType) (db DB, dbDir string) {
	dirname, err := ioutil.TempDir("", "db_common_test")
	require.NoError(t, err)
	db, err = NewDB("testdb", backend, dirname)
	require.NoError(t, err)
	return db, dirname
}

func newDB(t require.TestingT, backend BackendType) (DB, string, string) {
	name := fmt.Sprintf("test_%x", randStr(12))
	dir := os.TempDir()
	db, err := NewDB(name, backend, dir)
	require.NoError(t, err)
	return db, dir, name
}

func benchmarkRangeScans(b *testing.B, db DB, dbSize int64) {
	b.StopTimer()

	rangeSize := int64(10000)
	if dbSize < rangeSize {
		b.Errorf("db size %v cannot be less than range size %v", dbSize, rangeSize)
	}

	for i := int64(0); i < dbSize; i++ {
		bytes := int642Bytes(i)
		err := db.Set(bytes, bytes)
		if err != nil {
			// require.NoError() is very expensive (according to profiler), so check manually
			b.Fatal(b, err)
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		start := rand.Int63n(dbSize - rangeSize)
		end := start + rangeSize
		iter, err := db.Iterator(int642Bytes(start), int642Bytes(end))
		require.NoError(b, err)
		count := 0
		for ; iter.Valid(); iter.Next() {
			count++
		}
		iter.Close()
		require.EqualValues(b, rangeSize, count)
	}
}

func benchmarkRandomReadsWrites(b *testing.B, db DB) {
	const numItems = int64(1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRandomReadWrite(b, db, numItems)
	}
}

func benchmarkParallelRandomReadsWrites(b *testing.B, db DB) {
	const numItems = int64(1000000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkRandomReadWrite(b, db, numItems)
		}
	})
}

func benchmarkRandomReadWrite(b *testing.B, db DB, numItems int64) {
	// Write something
	{
		idx := rand.Int63n(numItems)
		val := idx + 2
		idxBytes := int642Bytes(idx)
		valBytes := int642Bytes(val)
		// fmt.Printf("Set %X -> %X\n", idxBytes, valBytes)
		err := db.Set(idxBytes, valBytes)
		if err != nil {
			// require.NoError() is very expensive (according to profiler), so check manually
			b.Fatal(b, err)
		}
	}

	// Read something
	{
		idx := rand.Int63n(numItems)
		valExp := idx + 2
		idxBytes := int642Bytes(idx)
		valBytes, err := db.Get(idxBytes)
		if err != nil {
			// require.NoError() is very expensive (according to profiler), so check manually
			b.Fatal(b, err)
		}
		// fmt.Printf("Get %X -> %X\n", idxBytes, valBytes)
		if valBytes != nil {
			if len(valBytes) != 8 {
				b.Errorf("Expected length 8 for %v, got %X", idx, valBytes)
				b.Fail()
			}
			valGot := bytes2Int64(valBytes)
			if valGot > 1 { // cleveldb return 1 even if `idx` doesn't exist
				if valExp != valGot {
					b.Errorf("Expected %v for %v, got %v", valExp, idx, valGot)
					b.Fail()
				}
			}
		}
	}
}

func int642Bytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func bytes2Int64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
