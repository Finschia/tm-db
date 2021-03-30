package metadb

import (
	"fmt"
	"testing"

	"github.com/line/tm-db/v2/internal/dbtest"
	"github.com/stretchr/testify/require"
)

// Empty iterator for empty db.
func TestPrefixIteratorNoMatchNil(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db, name, dir := newTempDB(t, backend)
			defer dbtest.CleanupDB(db, name, dir)

			itr, err := IteratePrefix(db, []byte("2"))
			require.NoError(t, err)
			defer itr.Close()

			dbtest.Invalid(t, itr)
		})
	}
}

// Empty iterator for db populated after iterator created.
func TestPrefixIteratorNoMatch1(t *testing.T) {
	for backend := range backends {
		if backend == BoltDBBackend {
			t.Log("bolt does not support concurrent writes while iterating")
			continue
		}

		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db, name, dir := newTempDB(t, backend)
			defer dbtest.CleanupDB(db, name, dir)

			itr, err := IteratePrefix(db, []byte("2"))
			require.NoError(t, err)
			defer itr.Close()

			err = db.SetSync([]byte("1"), []byte("value_1"))
			require.NoError(t, err)

			dbtest.Invalid(t, itr)
		})
	}
}

// Empty iterator for prefix starting after db entry.
func TestPrefixIteratorNoMatch2(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db, name, dir := newTempDB(t, backend)
			defer dbtest.CleanupDB(db, name, dir)

			err := db.SetSync([]byte("3"), []byte("value_3"))
			require.NoError(t, err)

			itr, err := IteratePrefix(db, []byte("4"))
			require.NoError(t, err)
			defer itr.Close()

			dbtest.Invalid(t, itr)
		})
	}
}

// Iterator with single val for db with single val, starting from that val.
func TestPrefixIteratorMatch1(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db, name, dir := newTempDB(t, backend)
			defer dbtest.CleanupDB(db, name, dir)

			err := db.SetSync([]byte("2"), []byte("value_2"))
			require.NoError(t, err)

			itr, err := IteratePrefix(db, []byte("2"))
			require.NoError(t, err)
			defer itr.Close()

			dbtest.Valid(t, itr, true)
			dbtest.Item(t, itr, []byte("2"), []byte("value_2"))
			dbtest.Next(t, itr, false)

			// Once invalid...
			dbtest.Invalid(t, itr)
		})
	}
}

// Iterator with prefix iterates over everything with same prefix.
func TestPrefixIteratorMatches1N(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db, name, dir := newTempDB(t, backend)
			defer dbtest.CleanupDB(db, name, dir)

			// prefixed
			err := db.SetSync([]byte("a/1"), []byte("value_1"))
			require.NoError(t, err)
			err = db.SetSync([]byte("a/3"), []byte("value_3"))
			require.NoError(t, err)

			// not
			err = db.SetSync([]byte("b/3"), []byte("value_3"))
			require.NoError(t, err)
			err = db.SetSync([]byte("a-3"), []byte("value_3"))
			require.NoError(t, err)
			err = db.SetSync([]byte("a.3"), []byte("value_3"))
			require.NoError(t, err)
			err = db.SetSync([]byte("abcdefg"), []byte("value_3"))
			require.NoError(t, err)
			itr, err := IteratePrefix(db, []byte("a/"))
			require.NoError(t, err)
			defer itr.Close()

			dbtest.Valid(t, itr, true)
			dbtest.Item(t, itr, []byte("a/1"), []byte("value_1"))
			dbtest.Next(t, itr, true)
			dbtest.Item(t, itr, []byte("a/3"), []byte("value_3"))

			// Bad!
			dbtest.Next(t, itr, false)

			// Once invalid...
			dbtest.Invalid(t, itr)
		})
	}
}
