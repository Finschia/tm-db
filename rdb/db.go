//go:build rocksdb
// +build rocksdb

package rdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"unsafe"

	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
)

type RDB struct {
	fn     string
	db     *C.rocksdb_t
	opts   *C.rocksdb_options_t
	ropts  *C.rocksdb_readoptions_t
	wopts  *C.rocksdb_writeoptions_t
	wsopts *C.rocksdb_writeoptions_t
}

type rdbIterator struct {
	db         *RDB
	it         *C.rocksdb_iterator_t
	ropts      *C.rocksdb_readoptions_t
	reverse    bool
	lowerBound []byte
	upperBound []byte
}

func cerror(cerr *C.char) error {
	if cerr == nil {
		return nil
	}
	err := errors.New(C.GoString(cerr))
	C.free(unsafe.Pointer(cerr))
	return err
}

func b2c(b []byte) *C.char {
	if len(b) == 0 {
		return nil
	}
	return (*C.char)(unsafe.Pointer(&b[0]))
}

func NewDB(name string, dir string) (*RDB, error) {
	var cerr *C.char

	fn := filepath.Join(dir, name+".db")

	bbto := C.rocksdb_block_based_options_create()
	C.rocksdb_block_based_options_set_block_cache(bbto, C.rocksdb_cache_create_lru(C.size_t(1<<30)))
	C.rocksdb_block_based_options_set_filter_policy(bbto, C.rocksdb_filterpolicy_create_bloom(C.int(10)))
	defer C.rocksdb_block_based_options_destroy(bbto)

	opts := C.rocksdb_options_create()
	C.rocksdb_options_set_block_based_table_factory(opts, bbto)
	C.rocksdb_options_set_create_if_missing(opts, 1)
	C.rocksdb_options_increase_parallelism(opts, C.int(runtime.NumCPU()))
	C.rocksdb_options_optimize_level_style_compaction(opts, 512*1024*1024)

	ropts := C.rocksdb_readoptions_create()
	wopts := C.rocksdb_writeoptions_create()
	wsopts := C.rocksdb_writeoptions_create()
	C.rocksdb_writeoptions_set_sync(wsopts, C.uchar(1))

	db := C.rocksdb_open(opts, b2c([]byte(fn)), &cerr)
	if cerr != nil {
		C.rocksdb_options_destroy(opts)
		C.rocksdb_writeoptions_destroy(wopts)
		C.rocksdb_writeoptions_destroy(wsopts)
		return nil, cerror(cerr)
	}
	return &RDB{
		fn:     fn,
		db:     db,
		opts:   opts,
		ropts:  ropts,
		wopts:  wopts,
		wsopts: wsopts,
	}, nil
}

func (db *RDB) Name() string {
	return db.fn
}

func (db *RDB) Get(key []byte) ([]byte, error) {
	var cerr *C.char
	var cvl C.size_t
	ck := b2c(key)
	cv := C.rocksdb_get(db.db, db.ropts, ck, C.size_t(len(key)), &cvl, &cerr)
	if cerr != nil {
		err := cerror(cerr)
		return nil, err
	}
	if cv == nil {
		return nil, nil
	}
	rv := C.GoBytes(unsafe.Pointer(cv), C.int(cvl))
	C.free(unsafe.Pointer(cv))
	return rv, nil
}

func (db *RDB) Has(key []byte) (bool, error) {
	var cerr *C.char
	var cvl C.size_t
	ck := b2c(key)
	cv := C.rocksdb_get(db.db, db.ropts, ck, C.size_t(len(key)), &cvl, &cerr)
	if cerr != nil {
		return false, cerror(cerr)
	}
	if cv == nil {
		return false, nil
	}
	C.free(unsafe.Pointer(cv))
	return true, nil
}

func (db *RDB) Set(key []byte, value []byte) error {
	var cerr *C.char
	ck, cv := b2c(key), b2c(value)
	C.rocksdb_put(db.db, db.wopts, ck, C.size_t(len(key)), cv, C.size_t(len(value)),
		&cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	return nil
}

func (db *RDB) SetSync(key []byte, value []byte) error {
	var cerr *C.char
	ck, cv := b2c(key), b2c(value)
	C.rocksdb_put(db.db, db.wsopts, ck, C.size_t(len(key)), cv, C.size_t(len(value)),
		&cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	return nil
}

func (db *RDB) Delete(key []byte) error {
	var cerr *C.char
	ck := b2c(key)
	C.rocksdb_delete(db.db, db.wopts, ck, C.size_t(len(key)), &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	return nil
}

func (db *RDB) DeleteSync(key []byte) error {
	var cerr *C.char
	ck := b2c(key)
	C.rocksdb_delete(db.db, db.wsopts, ck, C.size_t(len(key)), &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	return nil
}

func (db *RDB) Close() error {
	C.rocksdb_options_destroy(db.opts)
	C.rocksdb_readoptions_destroy(db.ropts)
	C.rocksdb_writeoptions_destroy(db.wopts)
	C.rocksdb_writeoptions_destroy(db.wsopts)
	C.rocksdb_close(db.db)
	return nil
}

// TODO: not implemented yet
func (db *RDB) Stats() map[string]string {
	m := map[string]string{}
	m["dummy"] = "100"
	return m
}

func (db *RDB) NewBatch() tmdb.Batch {
	return newRDBBatch(db)
}

func (db *RDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}
	return newRdbIterator(db, start, end, false), nil
}

func (db *RDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	return newRdbIterator(db, start, end, false), nil
}

func (db *RDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, tmdb.ErrKeyEmpty
	}
	return newRdbIterator(db, start, end, true), nil
}

func (db *RDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	return newRdbIterator(db, start, end, true), nil
}

func (db *RDB) Print() error {
	itr, err := db.Iterator(nil, nil)
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

// EOF
