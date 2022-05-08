//go:build rocksdb
// +build rocksdb

package db

import (
	"github.com/linxGnu/grocksdb"
)

type rocksDBIterator struct {
	source    *grocksdb.Iterator
	opts      *grocksdb.ReadOptions
	isReverse bool
	isInvalid bool
	key       []byte
	value     []byte
}

var _ Iterator = (*rocksDBIterator)(nil)

func newRockDBRangeOptions(start, end []byte) *grocksdb.ReadOptions {
	ro := grocksdb.NewDefaultReadOptions()
	if start != nil {
		ro.SetIterateLowerBound(start)
	}
	if end != nil {
		ro.SetIterateUpperBound(end)
	}
	return ro
}

func newRocksDBIterator(source *grocksdb.Iterator, opts *grocksdb.ReadOptions, isReverse bool) *rocksDBIterator {
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
		itr.key = moveSliceToBytes(itr.source.Key())
	}
	return itr.key
}

// Value implements Iterator.
func (itr *rocksDBIterator) Value() []byte {
	itr.assertIsValid()
	if itr.value == nil {
		itr.value = moveSliceToBytes(itr.source.Value())
	}
	return itr.value
}

// Next implements Iterator.
func (itr *rocksDBIterator) Next() {
	itr.assertIsValid()

	itr.key = nil
	itr.value = nil

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
	return nil
}

func (itr *rocksDBIterator) assertIsValid() {
	if itr.isInvalid {
		panic("iterator is invalid")
	}
}

// moveSliceToBytes will free the slice and copy out a go []byte
// This function can be applied on *Slice returned from Key() and Value()
// of an Iterator, because they are marked as freed.
func moveSliceToBytes(s *grocksdb.Slice) []byte {
	var bz []byte
	if s.Exists() {
		bz = cp(s.Data())
	}
	s.Free()
	return bz
}
