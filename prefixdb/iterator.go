package prefixdb

import (
	"bytes"
	tmdb "github.com/line/tm-db/v2"
)

// Strips prefix while iterating from Iterator.
type prefixDBIterator struct {
	prefix    []byte
	source    tmdb.Iterator
	isInvalid bool
}

var _ tmdb.Iterator = (*prefixDBIterator)(nil)

func newPrefixIterator(prefix []byte, source tmdb.Iterator) (*prefixDBIterator, error) {
	// Empty keys are not allowed, so if a key exists in the database that exactly matches the
	// prefix we need to skip it.
	if source.Valid() && bytes.Equal(source.Key(), prefix) {
		source.Next()
	}

	return &prefixDBIterator{
		prefix:    prefix,
		source:    source,
		isInvalid: !source.Valid(),
	}, nil
}

// Valid implements Iterator.
func (itr *prefixDBIterator) Valid() bool {
	// Once invalid, forever invalid.
	if itr.isInvalid {
		return false
	}

	// If source is invalid, invalid.
	if !itr.source.Valid() {
		itr.invalidate()
		return false
	}

	// Empty keys are not allowed, so if a key exists in the database that exactly matches the
	// prefix we need to skip it.
	if bytes.Equal(itr.source.Key(), itr.prefix) {
		itr.invalidate()
		return false
	}

	return true
}

func (itr *prefixDBIterator) invalidate() {
	itr.isInvalid = true
}

// Next implements Iterator.
func (itr *prefixDBIterator) Next() {
	itr.assertIsValid()
	itr.source.Next()
}

// Next implements Iterator.
func (itr *prefixDBIterator) Key() []byte {
	itr.assertIsValid()
	key := itr.source.Key()
	return key[len(itr.prefix):]
}

// Value implements Iterator.
func (itr *prefixDBIterator) Value() []byte {
	itr.assertIsValid()
	return itr.source.Value()
}

// Error implements Iterator.
func (itr *prefixDBIterator) Error() error {
	return itr.source.Error()
}

// Close implements Iterator.
func (itr *prefixDBIterator) Close() error {
	return itr.source.Close()
}

func (itr *prefixDBIterator) assertIsValid() {
	if itr.isInvalid {
		panic("iterator is invalid")
	}
}
