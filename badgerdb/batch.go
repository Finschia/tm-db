package badgerdb

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
	tmdb "github.com/line/tm-db/v2"
)

type badgerDBBatch struct {
	db *badger.DB
	wb *badger.WriteBatch

	// Calling db.Flush twice panics, so we must keep track of whether we've
	// flushed already on our own. If Write can receive from the firstFlush
	// channel, then it's the first and only Flush call we should do.
	//
	// Upstream bug report:
	// https://github.com/dgraph-io/badger/issues/1394
	firstFlush chan struct{}
}

var _ tmdb.Batch = (*badgerDBBatch)(nil)

func newBadgerDBBatch(b *BadgerDB) *badgerDBBatch {
	wb := &badgerDBBatch{
		db:         b.db,
		wb:         b.db.NewWriteBatch(),
		firstFlush: make(chan struct{}, 1),
	}
	wb.firstFlush <- struct{}{}
	return wb
}

func (b *badgerDBBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	return b.wb.Set(key, value)
}

func (b *badgerDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	return b.wb.Delete(key)
}

func (b *badgerDBBatch) Write() error {
	select {
	case <-b.firstFlush:
		return b.wb.Flush()
	default:
		return fmt.Errorf("batch already flushed")
	}
}

func (b *badgerDBBatch) WriteSync() error {
	return withSync(b.db, b.Write())
}

func (b *badgerDBBatch) Close() error {
	select {
	case <-b.firstFlush: // a Flush after Cancel panics too
	default:
	}
	b.wb.Cancel()
	return nil
}
