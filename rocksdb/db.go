package rocksdb

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/line/gorocksdb"
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

// RocksDB is a RocksDB backend.
type RocksDB struct {
	db     *gorocksdb.DB
	ro     *gorocksdb.ReadOptions
	wo     *gorocksdb.WriteOptions
	woSync *gorocksdb.WriteOptions
}

var _ tmdb.DB = (*RocksDB)(nil)

func NewDB(name string, dir string) (*RocksDB, error) {
	// default rocksdb option, good enough for most cases, including heavy workloads.
	// 1GB table cache, 512MB write buffer(may use 50% more on heavy workloads).
	// compression: snappy as default, need to -lsnappy to enable.
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(1 << 30))
	bbto.SetFilterPolicy(gorocksdb.NewBloomFilter(10))
	defer bbto.Destroy()

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.IncreaseParallelism(runtime.NumCPU())
	// 1.5GB maximum memory use for writebuffer.
	opts.OptimizeLevelStyleCompaction(512 * 1024 * 1024)
	defer opts.Destroy()

	return NewDBWithOptions(name, dir, opts)
}

func NewDBWithOptions(name string, dir string, opts *gorocksdb.Options) (*RocksDB, error) {
	err := util.MakePath(dir)
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, name+".db")
	db, err := gorocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil, err
	}
	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	woSync := gorocksdb.NewDefaultWriteOptions()
	woSync.SetSync(true)
	database := &RocksDB{
		db:     db,
		ro:     ro,
		wo:     wo,
		woSync: woSync,
	}
	return database, nil
}

// Get implements DB.
func (db *RocksDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}
	return db.db.GetBytes(db.ro, key)
}

// Has implements DB.
func (db *RocksDB) Has(key []byte) (bool, error) {
	s, err := db.db.Get(db.ro, key)
	if s != nil {
		defer s.Free()
	}
	if err != nil {
		return false, err
	}
	return s.Exists(), nil
}

// Set implements DB.
func (db *RocksDB) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return db.db.Put(db.wo, key, value)
}

// SetSync implements DB.
func (db *RocksDB) SetSync(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return db.db.Put(db.woSync, key, value)
}

// Delete implements DB.
func (db *RocksDB) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return db.db.Delete(db.wo, key)
}

// DeleteSync implements DB.
func (db *RocksDB) DeleteSync(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return db.db.Delete(db.woSync, key)
}

// FIXME This should not be exposed
func (db *RocksDB) DB() *gorocksdb.DB {
	return db.db
}

// Close implements DB.
func (db *RocksDB) Close() error {
	db.db.Close()
	db.ro.Destroy()
	db.wo.Destroy()
	db.woSync.Destroy()
	return nil
}

// Print implements DB.
func (db *RocksDB) Print() error {
	itr, err := db.Iterator(nil, nil)
	if err != nil {
		return err
	}
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		value := itr.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
	return nil
}

// Stats implements DB.
func (db *RocksDB) Stats() map[string]string {
	keys := []string{"rocksdb.stats"}
	stats := make(map[string]string, len(keys))
	for _, key := range keys {
		stats[key] = db.db.GetProperty(key)
	}
	return stats
}

// NewBatch implements DB.
func (db *RocksDB) NewBatch() tmdb.Batch {
	return newRocksDBBatch(db)
}

// Iterator implements DB.
func (db *RocksDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}

	ro := newRockDBRangeOptions(start, end)
	itr := db.db.NewIterator(ro)

	return newRocksDBIterator(itr, ro, false), nil
}

func (db *RocksDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}

	ro := newRockDBRangeOptions(start, end)
	itr := db.db.NewIterator(ro)

	return newRocksDBIterator(itr, ro, false), nil
}

// ReverseIterator implements DB.
func (db *RocksDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}

	ro := newRockDBRangeOptions(start, end)
	itr := db.db.NewIterator(ro)

	return newRocksDBIterator(itr, ro, true), nil
}

func (db *RocksDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}

	ro := newRockDBRangeOptions(start, end)
	itr := db.db.NewIterator(ro)

	return newRocksDBIterator(itr, ro, true), nil
}
