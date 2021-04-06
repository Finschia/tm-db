package dbtest

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	tmdb "github.com/line/tm-db/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------
// Helper functions.

func Valid(t *testing.T, itr tmdb.Iterator, expected bool) {
	valid := itr.Valid()
	require.Equal(t, expected, valid)
}

func Next(t *testing.T, itr tmdb.Iterator, expected bool) {
	itr.Next()
	// assert.NoError(t, err) TODO: look at fixing this
	valid := itr.Valid()
	require.Equal(t, expected, valid)
}

func NextPanics(t *testing.T, itr tmdb.Iterator) {
	assert.Panics(t, func() { itr.Next() }, "checkNextPanics expected an error but didn't")
}

func Item(t *testing.T, itr tmdb.Iterator, key []byte, value []byte) {
	v := itr.Value()

	k := itr.Key()

	assert.Exactly(t, key, k)
	assert.Exactly(t, value, v)
}

func Invalid(t *testing.T, itr tmdb.Iterator) {
	Valid(t, itr, false)
	KeyPanics(t, itr)
	ValuePanics(t, itr)
	NextPanics(t, itr)
}

func KeyPanics(t *testing.T, itr tmdb.Iterator) {
	assert.Panics(t, func() { itr.Key() }, "checkKeyPanics expected panic but didn't")
}

func Value(t *testing.T, db tmdb.DB, key []byte, valueWanted []byte) {
	valueGot, err := db.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, valueWanted, valueGot)
}

func ValuePanics(t *testing.T, itr tmdb.Iterator) {
	assert.Panics(t, func() { itr.Value() })
}

func NewTestName(prefix string) (name, dir string) {
	name = fmt.Sprintf("%s_%x", prefix, RandStr(12))
	dir = os.TempDir()
	return name, dir
}

func CleanupDB(db tmdb.DB, name, dir string) {
	if db != nil {
		_ = db.Close()
	}

	err := os.RemoveAll(filepath.Join(dir, name) + ".db")
	if err != nil {
		panic(err)
	}
}

func TestDBIterator(t *testing.T, db tmdb.DB) {
	for i := 0; i < 10; i++ {
		if i != 6 { // but skip 6.
			err := db.Set(Int642Bytes(int64(i)), []byte{})
			require.NoError(t, err)
		}
	}

	// Blank iterator keys should error
	_, err := db.Iterator([]byte{}, nil)
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.Iterator(nil, []byte{})
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.ReverseIterator([]byte{}, nil)
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.ReverseIterator(nil, []byte{})
	require.Equal(t, tmdb.ErrKeyEmpty, err)

	itr, err := db.Iterator(nil, nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{0, 1, 2, 3, 4, 5, 7, 8, 9}, "forward iterator")

	ritr, err := db.ReverseIterator(nil, nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{9, 8, 7, 5, 4, 3, 2, 1, 0}, "reverse iterator")

	itr, err = db.Iterator(nil, Int642Bytes(0))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64(nil), "forward iterator to 0")

	ritr, err = db.ReverseIterator(Int642Bytes(10), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64(nil), "reverse iterator from 10 (ex)")

	itr, err = db.Iterator(Int642Bytes(0), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{0, 1, 2, 3, 4, 5, 7, 8, 9}, "forward iterator from 0")

	itr, err = db.Iterator(Int642Bytes(1), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{1, 2, 3, 4, 5, 7, 8, 9}, "forward iterator from 1")

	ritr, err = db.ReverseIterator(nil, Int642Bytes(10))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{9, 8, 7, 5, 4, 3, 2, 1, 0}, "reverse iterator from 10 (ex)")

	ritr, err = db.ReverseIterator(nil, Int642Bytes(9))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{8, 7, 5, 4, 3, 2, 1, 0}, "reverse iterator from 9 (ex)")

	ritr, err = db.ReverseIterator(nil, Int642Bytes(8))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{7, 5, 4, 3, 2, 1, 0}, "reverse iterator from 8 (ex)")

	itr, err = db.Iterator(Int642Bytes(5), Int642Bytes(6))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{5}, "forward iterator from 5 to 6")

	itr, err = db.Iterator(Int642Bytes(5), Int642Bytes(7))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{5}, "forward iterator from 5 to 7")

	itr, err = db.Iterator(Int642Bytes(5), Int642Bytes(8))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{5, 7}, "forward iterator from 5 to 8")

	itr, err = db.Iterator(Int642Bytes(6), Int642Bytes(7))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64(nil), "forward iterator from 6 to 7")

	itr, err = db.Iterator(Int642Bytes(6), Int642Bytes(8))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{7}, "forward iterator from 6 to 8")

	itr, err = db.Iterator(Int642Bytes(7), Int642Bytes(8))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{7}, "forward iterator from 7 to 8")

	ritr, err = db.ReverseIterator(Int642Bytes(4), Int642Bytes(5))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{4}, "reverse iterator from 5 (ex) to 4")

	ritr, err = db.ReverseIterator(Int642Bytes(4), Int642Bytes(6))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{5, 4}, "reverse iterator from 6 (ex) to 4")

	ritr, err = db.ReverseIterator(Int642Bytes(4), Int642Bytes(7))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{5, 4}, "reverse iterator from 7 (ex) to 4")

	ritr, err = db.ReverseIterator(Int642Bytes(5), Int642Bytes(6))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{5}, "reverse iterator from 6 (ex) to 5")

	ritr, err = db.ReverseIterator(Int642Bytes(5), Int642Bytes(7))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{5}, "reverse iterator from 7 (ex) to 5")

	ritr, err = db.ReverseIterator(Int642Bytes(6), Int642Bytes(7))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64(nil), "reverse iterator from 7 (ex) to 6")

	ritr, err = db.ReverseIterator(Int642Bytes(10), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64(nil), "reverse iterator to 10")

	ritr, err = db.ReverseIterator(Int642Bytes(6), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{9, 8, 7}, "reverse iterator to 6")

	ritr, err = db.ReverseIterator(Int642Bytes(5), nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{9, 8, 7, 5}, "reverse iterator to 5")

	ritr, err = db.ReverseIterator(Int642Bytes(8), Int642Bytes(9))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{8}, "reverse iterator from 9 (ex) to 8")

	ritr, err = db.ReverseIterator(Int642Bytes(2), Int642Bytes(4))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64{3, 2}, "reverse iterator from 4 (ex) to 2")

	ritr, err = db.ReverseIterator(Int642Bytes(4), Int642Bytes(2))
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, []int64(nil), "reverse iterator from 2 (ex) to 4")
}

// Test `CONTRACT: No writes may happen within a domain while an iterator exists over it.`
func TestDBIteratorNoWrites(t *testing.T, db tmdb.DB) {
	for i := 0; i < 10; i++ {
		if i != 6 { // but skip 6.
			err := db.Set(Int642Bytes(int64(i)), []byte{})
			require.NoError(t, err)
		}
	}

	itr, err := db.Iterator(nil, nil)
	require.NoError(t, err)

	err = db.Set(Int642Bytes(int64(6)), []byte{})
	require.NoError(t, err)

	exist6, err := db.Has(Int642Bytes(int64(6)))
	require.True(t, exist6)

	verifyAndCloseIterator(t, itr, []int64{0, 1, 2, 3, 4, 5, 7, 8, 9}, "forward iterator")
}

func TestDBEmptyIterator(t *testing.T, db tmdb.DB) {
	itr, err := db.Iterator(nil, nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, nil, "forward iterator with empty db")

	ritr, err := db.ReverseIterator(nil, nil)
	require.NoError(t, err)
	verifyAndCloseIterator(t, ritr, nil, "reverse iterator with empty db")
}

func TestDBPrefixIterator(t *testing.T, db tmdb.DB) {
	for i := 0; i < 10; i++ {
		if i != 6 { // but skip 6.
			err := db.Set(Int642Bytes(int64(i)), []byte{})
			require.NoError(t, err)
		}
	}

	// Blank iterator keys should error
	_, err := db.PrefixIterator(nil)
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.PrefixIterator([]byte{})
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.ReversePrefixIterator(nil)
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	_, err = db.ReversePrefixIterator([]byte{})
	require.Equal(t, tmdb.ErrKeyEmpty, err)

	itr, err := db.PrefixIterator(Int642Bytes(0))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{0}, "forward iterator with 0 prefix")

	itr, err = db.ReversePrefixIterator(Int642Bytes(0))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, []int64{0}, "reverse iterator with 0 prefix")

	itr, err = db.PrefixIterator(Int642Bytes(6))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, nil, "forward iterator with 6 prefix")

	itr, err = db.ReversePrefixIterator(Int642Bytes(6))
	require.NoError(t, err)
	verifyAndCloseIterator(t, itr, nil, "reverse iterator with 6 prefix")
}

func verifyAndCloseIterator(t *testing.T, itr tmdb.Iterator, expected []int64, msg string) {
	var list []int64
	for itr.Valid() {
		key := itr.Key()
		list = append(list, Bytes2Int64(key))
		itr.Next()
	}
	assert.Equal(t, expected, list, msg)

	err := itr.Close()
	require.NoError(t, err)
}

// Empty iterator for empty db.
func TestPrefixIteratorNoMatchNil(t *testing.T, db tmdb.DB) {
	itr, err := db.PrefixIterator([]byte("2"))
	require.NoError(t, err)
	defer itr.Close()

	Invalid(t, itr)
}

// Empty iterator for db populated after iterator created.
func TestPrefixIteratorNoMatch1(t *testing.T, db tmdb.DB) {
	itr, err := db.PrefixIterator([]byte("2"))
	require.NoError(t, err)
	defer itr.Close()

	err = db.SetSync([]byte("1"), []byte("value_1"))
	require.NoError(t, err)

	Invalid(t, itr)
}

// Empty iterator for prefix starting after db entry.
func TestPrefixIteratorNoMatch2(t *testing.T, db tmdb.DB) {
	err := db.SetSync([]byte("3"), []byte("value_3"))
	require.NoError(t, err)

	itr, err := db.PrefixIterator([]byte("4"))
	require.NoError(t, err)
	defer itr.Close()

	Invalid(t, itr)
}

// Iterator with single val for db with single val, starting from that val.
func TestPrefixIteratorMatch1(t *testing.T, db tmdb.DB) {
	err := db.SetSync([]byte("2"), []byte("value_2"))
	require.NoError(t, err)

	itr, err := db.PrefixIterator([]byte("2"))
	require.NoError(t, err)
	defer itr.Close()

	Valid(t, itr, true)
	Item(t, itr, []byte("2"), []byte("value_2"))
	Next(t, itr, false)

	// Once invalid...
	Invalid(t, itr)
}

// Iterator with prefix iterates over everything with same prefix.
func TestPrefixIteratorMatches1N(t *testing.T, db tmdb.DB) {
	// prefixed
	err := db.SetSync([]byte("a/1"), []byte("value_1"))
	require.NoError(t, err)
	err = db.SetSync([]byte("a/3"), []byte("value_3"))
	require.NoError(t, err)

	// not
	err = db.SetSync([]byte("b/3"), []byte("value_3"))
	require.NoError(t, err)
	err = db.SetSync([]byte("a-3"), []byte("value_3"))
	require.NoError(t, err)
	err = db.SetSync([]byte("a.3"), []byte("value_3"))
	require.NoError(t, err)
	err = db.SetSync([]byte("abcdefg"), []byte("value_3"))
	require.NoError(t, err)
	itr, err := db.PrefixIterator([]byte("a/"))
	require.NoError(t, err)
	defer itr.Close()

	Valid(t, itr, true)
	Item(t, itr, []byte("a/1"), []byte("value_1"))
	Next(t, itr, true)
	Item(t, itr, []byte("a/3"), []byte("value_3"))

	// Bad!
	Next(t, itr, false)

	// Once invalid...
	Invalid(t, itr)
}

func TestDBBatch(t *testing.T, db tmdb.DB) {
	// create a new batch, and some items - they should not be visible until we write
	batch := db.NewBatch()
	require.NoError(t, batch.Set([]byte("a"), []byte{1}))
	require.NoError(t, batch.Set([]byte("b"), []byte{2}))
	require.NoError(t, batch.Set([]byte("c"), []byte{3}))
	assertKeyValues(t, db, map[string][]byte{})

	err := batch.Write()
	require.NoError(t, err)
	assertKeyValues(t, db, map[string][]byte{"a": {1}, "b": {2}, "c": {3}})

	// trying to modify or rewrite a written batch should error, but closing it should work
	require.Error(t, batch.Set([]byte("a"), []byte{9}))
	require.Error(t, batch.Delete([]byte("a")))
	require.Error(t, batch.Write())
	require.Error(t, batch.WriteSync())
	require.NoError(t, batch.Close())

	// batches should write changes in order
	batch = db.NewBatch()
	require.NoError(t, batch.Delete([]byte("a")))
	require.NoError(t, batch.Set([]byte("a"), []byte{1}))
	require.NoError(t, batch.Set([]byte("b"), []byte{1}))
	require.NoError(t, batch.Set([]byte("b"), []byte{2}))
	require.NoError(t, batch.Set([]byte("c"), []byte{3}))
	require.NoError(t, batch.Delete([]byte("c")))
	require.NoError(t, batch.Write())
	require.NoError(t, batch.Close())
	assertKeyValues(t, db, map[string][]byte{"a": {1}, "b": {2}})

	// empty and nil keys, as well as nil values, should be disallowed
	batch = db.NewBatch()
	err = batch.Set([]byte{}, []byte{0x01})
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	err = batch.Set(nil, []byte{0x01})
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	err = batch.Set([]byte("a"), nil)
	require.Equal(t, tmdb.ErrValueNil, err)

	err = batch.Delete([]byte{})
	require.Equal(t, tmdb.ErrKeyEmpty, err)
	err = batch.Delete(nil)
	require.Equal(t, tmdb.ErrKeyEmpty, err)

	err = batch.Close()
	require.NoError(t, err)

	// it should be possible to write an empty batch
	batch = db.NewBatch()
	err = batch.Write()
	require.NoError(t, err)
	assertKeyValues(t, db, map[string][]byte{"a": {1}, "b": {2}})

	// it should be possible to close an empty batch, and to re-close a closed batch
	batch = db.NewBatch()
	batch.Close()
	batch.Close()

	// all other operations on a closed batch should error
	require.Error(t, batch.Set([]byte("a"), []byte{9}))
	require.Error(t, batch.Delete([]byte("a")))
	require.Error(t, batch.Write())
	require.Error(t, batch.WriteSync())
}

func assertKeyValues(t *testing.T, db tmdb.DB, expect map[string][]byte) {
	iter, err := db.Iterator(nil, nil)
	require.NoError(t, err)
	defer iter.Close()

	actual := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		require.NoError(t, iter.Error())
		actual[string(iter.Key())] = iter.Value()
	}

	assert.Equal(t, expect, actual)
}

const strChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // 62 characters

// RandStr constructs a random alphanumeric string of given length.
func RandStr(length int) string {
	chars := []byte{}
MAIN_LOOP:
	for {
		val := rand.Int63() // nolint:gosec // G404: Use of weak random number generator
		for i := 0; i < 10; i++ {
			v := int(val & 0x3f) // rightmost 6 bits
			if v >= 62 {         // only 62 characters in strChars
				val >>= 6
				continue
			} else {
				chars = append(chars, strChars[v])
				if len(chars) == length {
					break MAIN_LOOP
				}
				val >>= 6
			}
		}
	}

	return string(chars)
}

func Int642Bytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func Bytes2Int64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func BenchmarkRangeScans(b *testing.B, db tmdb.DB, dbSize int64) {
	b.StopTimer()

	rangeSize := int64(10000)
	if dbSize < rangeSize {
		b.Errorf("db size %v cannot be less than range size %v", dbSize, rangeSize)
	}

	for i := int64(0); i < dbSize; i++ {
		bytes := Int642Bytes(i)
		err := db.Set(bytes, bytes)
		if err != nil {
			// require.NoError() is very expensive (according to profiler), so check manually
			b.Fatal(b, err)
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		start := rand.Int63n(dbSize - rangeSize) // nolint: gosec
		end := start + rangeSize
		iter, err := db.Iterator(Int642Bytes(start), Int642Bytes(end))
		require.NoError(b, err)
		count := 0
		for ; iter.Valid(); iter.Next() {
			count++
		}
		iter.Close()
		require.EqualValues(b, rangeSize, count)
	}
}

func BenchmarkRandomReadsWrites(b *testing.B, db tmdb.DB) {
	b.StopTimer()

	// create dummy data
	const numItems = int64(1000000)
	internal := map[int64]int64{}
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(0)
	}

	// fmt.Println("ok, starting")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Write something
		{
			idx := rand.Int63n(numItems) // nolint: gosec
			internal[idx]++
			val := internal[idx]
			idxBytes := Int642Bytes(idx)
			valBytes := Int642Bytes(val)
			// fmt.Printf("Set %X -> %X\n", idxBytes, valBytes)
			err := db.Set(idxBytes, valBytes)
			if err != nil {
				// require.NoError() is very expensive (according to profiler), so check manually
				b.Fatal(b, err)
			}
		}

		// Read something
		{
			idx := rand.Int63n(numItems) // nolint: gosec
			valExp := internal[idx]
			idxBytes := Int642Bytes(idx)
			valBytes, err := db.Get(idxBytes)
			if err != nil {
				// require.NoError() is very expensive (according to profiler), so check manually
				b.Fatal(b, err)
			}
			// fmt.Printf("Get %X -> %X\n", idxBytes, valBytes)
			if valExp == 0 {
				if !bytes.Equal(valBytes, nil) {
					b.Errorf("Expected %v for %v, got %X", nil, idx, valBytes)
					break
				}
			} else {
				if len(valBytes) != 8 {
					b.Errorf("Expected length 8 for %v, got %X", idx, valBytes)
					break
				}
				valGot := Bytes2Int64(valBytes)
				if valExp != valGot {
					b.Errorf("Expected %v for %v, got %v", valExp, idx, valGot)
					break
				}
			}
		}
	}
}
