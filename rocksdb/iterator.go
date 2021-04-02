package rocksdb

import (
	"bytes"

	tmdb "github.com/line/tm-db/v2"
	"github.com/tecbot/gorocksdb"
)

type rocksDBIterator struct {
	source     *gorocksdb.Iterator
	start, end []byte
	isReverse  bool
	isInvalid  bool
	key        []byte
	value      []byte
}

var _ tmdb.Iterator = (*rocksDBIterator)(nil)

func newRocksDBIterator(source *gorocksdb.Iterator, start, end []byte, isReverse bool) *rocksDBIterator {
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
				eoakey := moveSliceToBytes(source.Key()) // end or after key
				if bytes.Compare(end, eoakey) <= 0 {
					source.Prev()
				}
			} else {
				source.SeekToLast()
			}
		}
	}
	return &rocksDBIterator{
		source:    source,
		start:     start,
		end:       end,
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
	itr.source.Close()
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
func moveSliceToBytes(s *gorocksdb.Slice) []byte {
	defer s.Free()
	if !s.Exists() {
		return nil
	}
	v := make([]byte, len(s.Data()))
	copy(v, s.Data())
	return v
}
