//go:build cleveldb
// +build cleveldb

package db

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRandomReadsWrites2(b *testing.B) {
	b.StopTimer()

	numItems := int64(1000000)
	internal := map[int64]int64{}
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(0)
	}
	dir := os.TempDir()
	db, err := NewCLevelDB(fmt.Sprintf("test_%x", randStr(12)), dir)
	if err != nil {
		b.Fatal(err.Error())
		return
	}

	fmt.Println("ok, starting")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Write something
		{
			idx := (int64(rand.Int()) % numItems)
			internal[idx]++
			val := internal[idx]
			idxBytes := int642Bytes(idx)
			valBytes := int642Bytes(val)
			// fmt.Printf("Set %X -> %X\n", idxBytes, valBytes)
			// nolint: errcheck
			db.Set(
				idxBytes,
				valBytes,
			)
		}
		// Read something
		{
			idx := (int64(rand.Int()) % numItems)
			val := internal[idx]
			idxBytes := int642Bytes(idx)
			valBytes, err := db.Get(idxBytes)
			if err != nil {
				b.Error(err)
			}
			// fmt.Printf("Get %X -> %X\n", idxBytes, valBytes)
			if val == 0 {
				if !bytes.Equal(valBytes, nil) {
					b.Errorf("Expected %v for %v, got %X",
						nil, idx, valBytes)
					break
				}
			} else {
				if len(valBytes) != 8 {
					b.Errorf("Expected length 8 for %v, got %X",
						idx, valBytes)
					break
				}
				valGot := bytes2Int64(valBytes)
				if val != valGot {
					b.Errorf("Expected %v for %v, got %v",
						val, idx, valGot)
					break
				}
			}
		}
	}

	db.Close()
}

func TestCLevelDBNewDB(t *testing.T) {
	db, dir, name := newDB(t, CLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	_, ok := db.(*CLevelDB)
	assert.True(t, ok)
}

func TestCLevelDBStats(t *testing.T) {
	db, dir, name := newDB(t, CLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	assert.NotEmpty(t, db.Stats())
}

func BenchmarkCLevelDBRangeScans1M(b *testing.B) {
	db, dir, name := newDB(b, CLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(1e6))
}

func BenchmarkCLevelDBRangeScans10M(b *testing.B) {
	db, dir, name := newDB(b, CLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRangeScans(b, db, int64(10e6))
}

func BenchmarkCLevelDBRandomReadsWrites(b *testing.B) {
	db, dir, name := newDB(b, CLevelDBBackend)
	defer closeDBWithCleanupDBDir(db, dir, name)

	benchmarkRandomReadsWrites(b, db)
}
