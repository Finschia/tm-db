package prefixdb

import (
	"fmt"

	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

// PrefixDB wraps a namespace of another database as a logical database.
type PrefixDB struct {
	prefix []byte
	db     tmdb.DB
}

var _ tmdb.DB = (*PrefixDB)(nil)

// NewPrefixDB lets you namespace multiple DBs within a single DB.
func NewDB(db tmdb.DB, prefix []byte) *PrefixDB {
	// TODO return error
	if len(prefix) == 0 {
		panic("prefix should not be empty")
	}
	return &PrefixDB{
		prefix: prefix,
		db:     db,
	}
}

func (pdb *PrefixDB) Name() string {
	return fmt.Sprintf("prefix(%s, %s)", pdb.db.Name(), pdb.prefix)
}

// Get implements DB.
func (pdb *PrefixDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}
	return pdb.db.Get(pdb.prefixed(key))
}

// Has implements DB.
func (pdb *PrefixDB) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, tmdb.ErrKeyEmpty
	}
	return pdb.db.Has(pdb.prefixed(key))
}

// Set implements DB.
func (pdb *PrefixDB) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return pdb.db.Set(pdb.prefixed(key), value)
}

// SetSync implements DB.
func (pdb *PrefixDB) SetSync(key []byte, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return pdb.db.SetSync(pdb.prefixed(key), value)
}

// Delete implements DB.
func (pdb *PrefixDB) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return pdb.db.Delete(pdb.prefixed(key))
}

// DeleteSync implements DB.
func (pdb *PrefixDB) DeleteSync(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return pdb.db.DeleteSync(pdb.prefixed(key))
}

// Iterator implements DB.
func (pdb *PrefixDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}

	var itr tmdb.Iterator
	var err error
	if start == nil && end == nil {
		// NOTE `PrefixIterator` has more chance to optimize performance
		itr, err = pdb.db.PrefixIterator(pdb.prefix)
	} else {
		pstart, pend := pdb.prefixedRange(start, end)
		itr, err = pdb.db.Iterator(pstart, pend)
	}
	if err != nil {
		return nil, err
	}

	return newPrefixIterator(pdb.prefix, itr)
}

func (pdb *PrefixDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	if len(prefix) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}

	itr, err := pdb.db.PrefixIterator(pdb.prefixed(prefix))
	if err != nil {
		return nil, err
	}

	return newPrefixIterator(pdb.prefix, itr)
}

// ReverseIterator implements DB.
func (pdb *PrefixDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}

	var ritr tmdb.Iterator
	var err error
	if start == nil && end == nil {
		// NOTE `ReversePrefixIterator` has more chance to optimize performance
		ritr, err = pdb.db.ReversePrefixIterator(pdb.prefix)
	} else {
		pstart, pend := pdb.prefixedRange(start, end)
		ritr, err = pdb.db.ReverseIterator(pstart, pend)
	}
	if err != nil {
		return nil, err
	}

	return newPrefixIterator(pdb.prefix, ritr)
}

func (pdb *PrefixDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	if len(prefix) == 0 {
		return nil, tmdb.ErrKeyEmpty
	}

	itr, err := pdb.db.ReversePrefixIterator(pdb.prefixed(prefix))
	if err != nil {
		return nil, err
	}

	return newPrefixIterator(pdb.prefix, itr)
}

// NewBatch implements DB.
func (pdb *PrefixDB) NewBatch() tmdb.Batch {
	return newPrefixBatch(pdb.prefix, pdb.db.NewBatch())
}

// Close implements DB.
func (pdb *PrefixDB) Close() error {
	return pdb.db.Close()
}

// Print implements DB.
func (pdb *PrefixDB) Print() error {
	fmt.Printf("prefix: %X\n", pdb.prefix)

	itr, err := pdb.Iterator(nil, nil)
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
func (pdb *PrefixDB) Stats() map[string]string {
	stats := make(map[string]string)
	stats["prefixdb.prefix.string"] = string(pdb.prefix)
	stats["prefixdb.prefix.hex"] = fmt.Sprintf("%X", pdb.prefix)
	source := pdb.db.Stats()
	for key, value := range source {
		stats["prefixdb.source."+key] = value
	}
	return stats
}

func (pdb *PrefixDB) prefixed(key []byte) []byte {
	return util.Concat(pdb.prefix, key)
}

func (pdb *PrefixDB) prefixedRange(start, end []byte) (pstart, pend []byte) {
	pstart = util.Concat(pdb.prefix, start)
	if end == nil {
		pend = util.CpIncr(pdb.prefix)
	} else {
		pend = util.Concat(pdb.prefix, end)
	}
	return pstart, pend
}
