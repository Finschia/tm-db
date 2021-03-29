package prefixdb

import (
	tmdb "github.com/line/tm-db/v2"
)

type prefixDBBatch struct {
	prefix []byte
	source tmdb.Batch
}

var _ tmdb.Batch = (*prefixDBBatch)(nil)

func newPrefixBatch(prefix []byte, source tmdb.Batch) prefixDBBatch {
	return prefixDBBatch{
		prefix: prefix,
		source: source,
	}
}

// Set implements Batch.
func (pb prefixDBBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	if value == nil {
		return tmdb.ErrValueNil
	}
	pkey := tmdb.Concat(pb.prefix, key)
	return pb.source.Set(pkey, value)
}

// Delete implements Batch.
func (pb prefixDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return tmdb.ErrKeyEmpty
	}
	pkey := tmdb.Concat(pb.prefix, key)
	return pb.source.Delete(pkey)
}

// Write implements Batch.
func (pb prefixDBBatch) Write() error {
	return pb.source.Write()
}

// WriteSync implements Batch.
func (pb prefixDBBatch) WriteSync() error {
	return pb.source.WriteSync()
}

// Close implements Batch.
func (pb prefixDBBatch) Close() error {
	return pb.source.Close()
}
