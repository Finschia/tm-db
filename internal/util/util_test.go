package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkConcat(b *testing.B) {
	bz1 := []byte("prefix")
	bz2 := []byte("key")
	for i := 0; i < b.N; i++ {
		_ = Concat(bz1, bz2)
	}
}

func BenchmarkPrefixed(b *testing.B) {
	bz1 := []byte("prefix")
	bz2 := []byte("key")
	for i := 0; i < b.N; i++ {
		_ = append(Cp(bz1), bz2...)
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
	require.Equal(t, bytes.Join([][]byte{prefix, key}, nil), Concat(prefix, key))
	require.Equal(t, prefix, Concat(prefix, nil))
	require.Equal(t, key, Concat(nil, key))
	require.Equal(t, []byte{}, Concat(nil, nil))
}
