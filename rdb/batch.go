//go:build rocksdb
// +build rocksdb

package rdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	tmdb "github.com/line/tm-db/v2"
)

type rdbBatch struct {
	db *RDB
	b  *C.rocksdb_writebatch_t
}

func newRDBBatch(db *RDB) *rdbBatch {
	return &rdbBatch{
		db: db,
		b:  C.rocksdb_writebatch_create(),
	}
}

func (b *rdbBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	if b.b == nil {
		return tmdb.ErrBatchClosed
	}
	ck, cv := b2c(key), b2c(value)
	C.rocksdb_writebatch_put(b.b, ck, C.size_t(len(key)), cv, C.size_t(len(value)))
	return nil
}

func (b *rdbBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if b.b == nil {
		return tmdb.ErrBatchClosed
	}
	C.rocksdb_writebatch_delete(b.b, b2c(key), C.size_t(len(key)))
	return nil
}

func (b *rdbBatch) Write() error {
	if b.b == nil {
		return tmdb.ErrBatchClosed
	}
	var cerr *C.char
	C.rocksdb_write(b.db.db, b.db.wopts, b.b, &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	b.Close()
	return nil
}

func (b *rdbBatch) WriteSync() error {
	if b.b == nil {
		return tmdb.ErrBatchClosed
	}
	var cerr *C.char
	C.rocksdb_write(b.db.db, b.db.wsopts, b.b, &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	b.Close()
	return nil
}

func (b *rdbBatch) WriteLowPri() error {
	if b.b == nil {
		return tmdb.ErrBatchClosed
	}
	var cerr *C.char
	C.rocksdb_write(b.db.db, b.db.wlpopts, b.b, &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	b.Close()
	return nil
}

func (b *rdbBatch) Close() error {
	if b.b != nil {
		C.rocksdb_writebatch_destroy(b.b)
		b.b = nil
	}
	return nil
}
