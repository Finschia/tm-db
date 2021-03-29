package db

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkConcat(b *testing.B) {
	bz1 := []byte("prefix")
	bz2 := []byte("key")
	for i := 0; i < b.N; i++ {
		_ = concat(bz1, bz2)
	}
}

func BenchmarkPrefixed(b *testing.B) {
	bz1 := []byte("prefix")
	bz2 := []byte("key")
	for i := 0; i < b.N; i++ {
		_ = append(cp(bz1), bz2...)
	}
}

func BenchmarkBytesJoin(b *testing.B) {
	bzz := [][]byte{[]byte("prefix"), []byte("key")}
	for i := 0; i < b.N; i++ {
		_ = bytes.Join(bzz, nil)
	}
}

func TestConcat(t *testing.T) {
	prefix := []byte("prefix")
	key := []byte("key")
	require.Equal(t, bytes.Join([][]byte{prefix, key}, nil), concat(prefix, key))
	require.Equal(t, prefix, concat(prefix, nil))
	require.Equal(t, key, concat(nil, key))
	require.Equal(t, []byte{}, concat(nil, nil))
}
