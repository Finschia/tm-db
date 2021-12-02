//go:build rocksdb
// +build rocksdb

package rdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"unsafe"
)

func newRdbIterator(db *RDB, lowerBound, upperBound []byte, reverse bool) *rdbIterator {
	ropts := C.rocksdb_readoptions_create()
	if len(lowerBound) >= 0 {
		C.rocksdb_readoptions_set_iterate_lower_bound(ropts, b2c(lowerBound), C.size_t(len(lowerBound)))
	}
	if len(upperBound) >= 0 {
		C.rocksdb_readoptions_set_iterate_upper_bound(ropts, b2c(upperBound), C.size_t(len(upperBound)))
	}
	it := C.rocksdb_create_iterator(db.db, ropts)
	if !reverse {
		C.rocksdb_iter_seek_to_first(it)
	} else {
		C.rocksdb_iter_seek_to_last(it)
	}
	return &rdbIterator{
		db:         db,
		it:         it,
		ropts:      ropts,
		reverse:    reverse,
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
}

func (itr *rdbIterator) Valid() bool {
	return C.rocksdb_iter_valid(itr.it) != 0
}

func (itr *rdbIterator) Key() []byte {
	if C.rocksdb_iter_valid(itr.it) == 0 {
		panic("iterator is invalid")
	}

	var cvl C.size_t
	cv := C.rocksdb_iter_key(itr.it, &cvl)
	if cv == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(cv), C.int(cvl))
}

func (itr *rdbIterator) Value() []byte {
	if C.rocksdb_iter_valid(itr.it) == 0 {
		panic("iterator is invalid")
	}
	var cvl C.size_t
	cv := C.rocksdb_iter_value(itr.it, &cvl)
	if cv == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(cv), C.int(cvl))
}

func (itr *rdbIterator) Next() {
	if C.rocksdb_iter_valid(itr.it) == 0 {
		panic("iterator is invalid")
	}
	if !itr.reverse {
		C.rocksdb_iter_next(itr.it)
	} else {
		C.rocksdb_iter_prev(itr.it)
	}
}

func (itr *rdbIterator) Error() error {
	var cerr *C.char
	if itr.it == nil {
		return nil
	}
	C.rocksdb_iter_get_error(itr.it, &cerr)
	if cerr != nil {
		return cerror(cerr)
	}
	return nil
}

func (itr *rdbIterator) Close() error {
	if itr.it != nil {
		C.rocksdb_iter_destroy(itr.it)
	}
	itr.it, itr.db, itr.lowerBound, itr.upperBound = nil, nil, nil, nil
	return nil
}

// EOF
