//go:build cleveldb
// +build cleveldb

package db

import (
	"bytes"

	"github.com/jmhodges/levigo"
)

// cLevelDBIterator is a cLevelDB iterator.
type cLevelDBIterator struct {
	source     *levigo.Iterator
	start, end []byte
	isReverse  bool
	isInvalid  bool
	key        []byte
	value      []byte
}

var _ Iterator = (*cLevelDBIterator)(nil)

func newCLevelDBIterator(source *levigo.Iterator, start, end []byte, isReverse bool) *cLevelDBIterator {
	if !isReverse {
		if len(start) == 0 {
			source.SeekToFirst()
		} else {
			source.Seek(start)
		}
	} else {
		if len(end) == 0 {
			source.SeekToLast()
		} else {
			source.Seek(end)
			if source.Valid() {
				eoakey := source.Key() // end or after key
				if bytes.Compare(end, eoakey) <= 0 {
					source.Prev()
				}
			} else {
				source.SeekToLast()
			}
		}
	}
	return &cLevelDBIterator{
		source:    source,
		start:     start,
		end:       end,
		isReverse: isReverse,
		isInvalid: false,
	}
}

// Valid implements Iterator.
func (itr *cLevelDBIterator) Valid() bool {
	// Once invalid, forever invalid.
	if itr.isInvalid {
		return false
	}

	// If source is invalid, invalid.
	if !itr.source.Valid() {
		itr.invalidate()
		return false
	}

	// If key is end or past it, invalid.
	var start = itr.start
	var end = itr.end
	var key = itr.Key()
	if !itr.isReverse {
		if end != nil && bytes.Compare(end, key) <= 0 {
			itr.invalidate()
			return false
		}
	} else {
		if start != nil && bytes.Compare(key, start) < 0 {
			itr.invalidate()
			return false
		}
	}

	// It's valid.
	return true
}

func (itr *cLevelDBIterator) invalidate() {
	itr.isInvalid = true
	itr.key = nil
	itr.value = nil
}

// Key implements Iterator.
func (itr *cLevelDBIterator) Key() []byte {
	itr.assertIsValid()
	if itr.key == nil {
		itr.key = itr.source.Key()
	}
	return itr.key
}

// Value implements Iterator.
func (itr *cLevelDBIterator) Value() []byte {
	itr.assertIsValid()
	if itr.value == nil {
		itr.value = itr.source.Value()
	}
	return itr.value
}

// Next implements Iterator.
func (itr *cLevelDBIterator) Next() {
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
func (itr *cLevelDBIterator) Error() error {
	return itr.source.GetError()
}

// Close implements Iterator.
func (itr *cLevelDBIterator) Close() error {
	itr.source.Close()
	return nil
}

func (itr *cLevelDBIterator) assertIsValid() {
	if itr.isInvalid {
		panic("iterator is invalid")
	}
}
