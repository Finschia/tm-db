package cleveldb

import (
	"fmt"
	"path/filepath"

	"github.com/jmhodges/levigo"
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

// CLevelDB uses the C LevelDB database via a Go wrapper.
type CLevelDB struct {
	name   string
	db     *levigo.DB
	ro     *levigo.ReadOptions
	wo     *levigo.WriteOptions
	woSync *levigo.WriteOptions
}

var _ tmdb.DB = (*CLevelDB)(nil)

// New creates a new CLevelDB.
func NewDB(name string, dir string) (*CLevelDB, error) {
	err := util.MakePath(dir)
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, name+".db")

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(1 << 30))
	opts.SetCreateIfMissing(true)
	defer opts.Close()

	db, err := levigo.Open(dbPath, opts)
	if err != nil {
		return nil, err
	}
	ro := levigo.NewReadOptions()
	wo := levigo.NewWriteOptions()
	woSync := levigo.NewWriteOptions()
	woSync.SetSync(true)
	database := &CLevelDB{
		name:   dbPath,
		db:     db,
		ro:     ro,
		wo:     wo,
		woSync: woSync,
	}
	return database, nil
}

func (db *CLevelDB) Name() string {
	return db.name
}

// Get implements DB.
func (db *CLevelDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}
	return db.db.Get(db.ro, key)
}

// Has implements DB.
func (db *CLevelDB) Has(key []byte) (bool, error) {
	bytes, err := db.Get(key)
	if err != nil {
		return false, err
	}
	return bytes != nil, nil
}

// Set implements DB.
func (db *CLevelDB) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return db.db.Put(db.wo, key, value)
}

// SetSync implements DB.
func (db *CLevelDB) SetSync(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return db.db.Put(db.woSync, key, value)
}

// Delete implements DB.
func (db *CLevelDB) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return db.db.Delete(db.wo, key)
}

// DeleteSync implements DB.
func (db *CLevelDB) DeleteSync(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return db.db.Delete(db.woSync, key)
}

// FIXME This should not be exposed
func (db *CLevelDB) DB() *levigo.DB {
	return db.db
}

// Close implements DB.
func (db *CLevelDB) Close() error {
	db.db.Close()
	db.ro.Close()
	db.wo.Close()
	db.woSync.Close()
	return nil
}

// Print implements DB.
func (db *CLevelDB) Print() error {
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
func (db *CLevelDB) Stats() map[string]string {
	keys := []string{
		"leveldb.aliveiters",
		"leveldb.alivesnaps",
		"leveldb.blockpool",
		"leveldb.cachedblock",
		"leveldb.num-files-at-level{n}",
		"leveldb.openedtables",
		"leveldb.sstables",
		"leveldb.stats",
	}

	stats := make(map[string]string, len(keys))
	for _, key := range keys {
		str := db.db.PropertyValue(key)
		stats[key] = str
	}
	return stats
}

// NewBatch implements DB.
func (db *CLevelDB) NewBatch() tmdb.Batch {
	return newCLevelDBBatch(db)
}

// Iterator implements DB.
func (db *CLevelDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}
	itr := db.db.NewIterator(db.ro)
	return newCLevelDBIterator(itr, start, end, false), nil
}

func (db *CLevelDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	itr := db.db.NewIterator(db.ro)
	return newCLevelDBIterator(itr, start, end, false), nil
}

// ReverseIterator implements DB.
func (db *CLevelDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}
	itr := db.db.NewIterator(db.ro)
	return newCLevelDBIterator(itr, start, end, true), nil
}

func (db *CLevelDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	itr := db.db.NewIterator(db.ro)
	return newCLevelDBIterator(itr, start, end, true), nil
}
