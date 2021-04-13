package rocksdb

import (
	"github.com/line/gorocksdb"
	tmdb "github.com/line/tm-db/v2"
)

type rocksDBIterator struct {
	source    *gorocksdb.Iterator
	opts      *gorocksdb.ReadOptions
	isReverse bool
	isInvalid bool
	key       *gorocksdb.Slice
	value     *gorocksdb.Slice
}

var _ tmdb.Iterator = (*rocksDBIterator)(nil)

func newRockDBRangeOptions(start, end []byte) *gorocksdb.ReadOptions {
	ro := gorocksdb.NewDefaultReadOptions()
	if start != nil {
		ro.SetIterateLowerBound(start)
	}
	if end != nil {
		ro.SetIterateUpperBound(end)
	}
	return ro
}

func newRocksDBIterator(source *gorocksdb.Iterator, opts *gorocksdb.ReadOptions, isReverse bool) *rocksDBIterator {
	if !isReverse {
		source.SeekToFirst()
	} else {
		source.SeekToLast()
	}

	return &rocksDBIterator{
		source:    source,
		opts:      opts,
		isReverse: isReverse,
		isInvalid: false,
	}
}

// Valid implements Iterator.
func (itr *rocksDBIterator) Valid() bool {
	// Once invalid, forever invalid.
	if itr.isInvalid {
		return false
	}

	// If source is invalid, invalid.
	if !itr.source.Valid() {
		itr.invalidate()
		return false
	}

	// It's valid.
	return true
}

func (itr *rocksDBIterator) invalidate() {
	itr.isInvalid = true
	itr.key = nil
	itr.value = nil
}

// Key implements Iterator.
func (itr *rocksDBIterator) Key() []byte {
	itr.assertIsValid()
	if itr.key == nil {
		itr.key = itr.source.Key()
	}
	return itr.key.Data()
}

// Value implements Iterator.
func (itr *rocksDBIterator) Value() []byte {
	itr.assertIsValid()
	if itr.value == nil {
		itr.value = itr.source.Value()
	}
	return itr.value.Data()
}

// Next implements Iterator.
func (itr *rocksDBIterator) Next() {
	itr.assertIsValid()

	itr.freeKeyValue()

	if !itr.isReverse {
		itr.source.Next()
	} else {
		itr.source.Prev()
	}
}

// Error implements Iterator.
func (itr *rocksDBIterator) Error() error {
	return itr.source.Err()
}

// Close implements Iterator.
func (itr *rocksDBIterator) Close() error {
	if itr.source != nil {
		itr.source.Close()
		itr.source = nil
	}
	if itr.opts != nil {
		itr.opts.Destroy()
		itr.opts = nil
	}
	itr.freeKeyValue()
	return nil
}

func (itr *rocksDBIterator) freeKeyValue() {
	if itr.key != nil {
		itr.key.Free()
		itr.key = nil
	}
	if itr.value != nil {
		itr.value.Free()
		itr.value = nil
	}
}

func (itr *rocksDBIterator) assertIsValid() {
	if itr.isInvalid {
		panic("iterator is invalid")
	}
}
