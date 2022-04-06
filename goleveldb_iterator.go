package db

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type goLevelDBIterator struct {
	source    iterator.Iterator
	isReverse bool
}

var _ Iterator = (*goLevelDBIterator)(nil)

func newGoLevelDBIterator(source iterator.Iterator, isReverse bool) *goLevelDBIterator {
	if !isReverse {
		source.First()
	} else {
		source.Last()
	}

	return &goLevelDBIterator{
		source:    source,
		isReverse: isReverse,
	}
}

// Valid implements Iterator.
func (itr *goLevelDBIterator) Valid() bool {
	return itr.source.Valid()
}

// Key implements Iterator.
func (itr *goLevelDBIterator) Key() []byte {
	// Key returns a copy of the current key.
	// See https://github.com/syndtr/goleveldb/blob/52c212e6c196a1404ea59592d3f1c227c9f034b2/leveldb/iterator/iter.go#L88
	itr.assertIsValid()
	return cp(itr.source.Key())
}

// Value implements Iterator.
func (itr *goLevelDBIterator) Value() []byte {
	// Value returns a copy of the current value.
	// See https://github.com/syndtr/goleveldb/blob/52c212e6c196a1404ea59592d3f1c227c9f034b2/leveldb/iterator/iter.go#L88
	itr.assertIsValid()
	return cp(itr.source.Value())
}

// Next implements Iterator.
func (itr *goLevelDBIterator) Next() {
	itr.assertIsValid()
	if !itr.isReverse {
		itr.source.Next()
	} else {
		itr.source.Prev()
	}
}

// Error implements Iterator.
func (itr *goLevelDBIterator) Error() error {
	return itr.source.Error()
}

// Close implements Iterator.
func (itr *goLevelDBIterator) Close() error {
	itr.source.Release()
	return nil
}

func (itr *goLevelDBIterator) assertIsValid() {
	if !itr.Valid() {
		panic("iterator is invalid")
	}
}
