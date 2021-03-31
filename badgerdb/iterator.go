package badgerdb

import (
	"bytes"

	"github.com/dgraph-io/badger/v2"
	tmdb "github.com/line/tm-db/v2"
)

type badgerDBIterator struct {
	reverse    bool
	start, end []byte

	txn  *badger.Txn
	iter *badger.Iterator

	lastErr error
}

var _ tmdb.Iterator = (*badgerDBIterator)(nil)

func newBadgerDBIterator(b *BadgerDB, start, end []byte, opts badger.IteratorOptions) (*badgerDBIterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}
	txn := b.db.NewTransaction(false)
	iter := txn.NewIterator(opts)
	iter.Rewind()
	iter.Seek(start)
	if opts.Reverse && iter.Valid() && bytes.Equal(iter.Item().Key(), start) {
		// If we're going in reverse, our starting point was "end",
		// which is exclusive.
		iter.Next()
	}
	return &badgerDBIterator{
		reverse: opts.Reverse,
		start:   start,
		end:     end,

		txn:  txn,
		iter: iter,
	}, nil
}

func (i *badgerDBIterator) Close() error {
	i.iter.Close()
	i.txn.Discard()
	return nil
}

func (i *badgerDBIterator) Domain() (start, end []byte) {
	return i.start, i.end
}

func (i *badgerDBIterator) Error() error {
	return i.lastErr
}

func (i *badgerDBIterator) Next() {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	i.iter.Next()
}

func (i *badgerDBIterator) Valid() bool {
	if !i.iter.Valid() {
		return false
	}
	if len(i.end) > 0 {
		key := i.iter.Item().Key()
		if c := bytes.Compare(key, i.end); (!i.reverse && c >= 0) || (i.reverse && c < 0) {
			// We're at the end key, or past the end.
			return false
		}
	}
	return true
}

func (i *badgerDBIterator) Key() []byte {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	// Note that we don't use KeyCopy, so this is only valid until the next
	// call to Next.
	return i.iter.Item().KeyCopy(nil)
}

func (i *badgerDBIterator) Value() []byte {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	val, err := i.iter.Item().ValueCopy(nil)
	if err != nil {
		i.lastErr = err
	}
	return val
}
