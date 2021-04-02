package badgerdb

import (
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

type BadgerDB struct {
	db *badger.DB
}

var _ tmdb.DB = (*BadgerDB)(nil)

// NewDB creates a Badger key-value store backed to the
// directory dir supplied. If dir does not exist, it will be created.
func NewDB(dbName, dir string) (*BadgerDB, error) {
	// Since Badger doesn't support database names, we join both to obtain
	// the final directory to use for the database.
	path := filepath.Join(dir, dbName)

	if err := util.MakePath(path); err != nil {
		return nil, err
	}
	opts := badger.DefaultOptions(path)
	opts.SyncWrites = false // note that we have Sync methods
	opts.Logger = nil       // badger is too chatty by default
	return NewDBWithOptions(opts)
}

// NewDBWithOptions creates a BadgerDB key value store
// gives the flexibility of initializing a database with the
// respective options.
func NewDBWithOptions(opts badger.Options) (*BadgerDB, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{db: db}, nil
}

func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err == nil && val == nil {
			val = []byte{}
		}
		return err
	})
	return val, err
}

func (b *BadgerDB) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, tmdb.ErrKeyEmpty
	}
	var found bool
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		found = err != badger.ErrKeyNotFound
		return nil
	})
	return found, err
}

func (b *BadgerDB) Set(key, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func withSync(db *badger.DB, err error) error {
	if err != nil {
		return err
	}
	return db.Sync()
}

func (b *BadgerDB) SetSync(key, value []byte) error {
	return withSync(b.db, b.Set(key, value))
}

func (b *BadgerDB) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerDB) DeleteSync(key []byte) error {
	return withSync(b.db, b.Delete(key))
}

func (b *BadgerDB) Close() error {
	return b.db.Close()
}

func (b *BadgerDB) Print() error {
	return nil
}

func (b *BadgerDB) Stats() map[string]string {
	return nil
}

func (b *BadgerDB) NewBatch() tmdb.Batch {
	return newBadgerDBBatch(b)
}

func (b *BadgerDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	opts := badger.DefaultIteratorOptions
	return newBadgerDBIterator(b, start, end, opts)
}

func (b *BadgerDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixRange(prefix)
	if err != nil {
		return nil, err
	}
	return b.Iterator(start, end)
}

func (b *BadgerDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	opts := badger.DefaultIteratorOptions
	opts.Reverse = true
	return newBadgerDBIterator(b, end, start, opts)
}

func (b *BadgerDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixRange(prefix)
	if err != nil {
		return nil, err
	}
	return b.ReverseIterator(start, end)
}
