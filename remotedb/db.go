package remotedb

import (
	"context"
	"errors"
	"fmt"

	tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/internal/util"
	"github.com/line/tm-db/v2/remotedb/grpcdb"
	protodb "github.com/line/tm-db/v2/remotedb/proto"
)

type RemoteDB struct {
	ctx context.Context
	dc  protodb.DBClient
}

func NewDB(serverAddr string, serverKey string) (*RemoteDB, error) {
	return newDB(grpcdb.NewClient(serverAddr, serverKey))
}

func newDB(gdc protodb.DBClient, err error) (*RemoteDB, error) {
	if err != nil {
		return nil, err
	}
	return &RemoteDB{dc: gdc, ctx: context.Background()}, nil
}

type Init struct {
	Dir  string
	Name string
	Type string
}

func (rd *RemoteDB) Name() string {
	return "remote"
}

func (rd *RemoteDB) InitRemote(in *Init) error {
	_, err := rd.dc.Init(rd.ctx, &protodb.Init{Dir: in.Dir, Type: in.Type, Name: in.Name})
	return err
}

var _ tmdb.DB = (*RemoteDB)(nil)

// Close is a noop currently
func (rd *RemoteDB) Close() error {
	return nil
}

func (rd *RemoteDB) Delete(key []byte) error {
	if _, err := rd.dc.Delete(rd.ctx, &protodb.Entity{Key: key}); err != nil {
		return fmt.Errorf("remoteDB.Delete: %w", err)
	}
	return nil
}

func (rd *RemoteDB) DeleteSync(key []byte) error {
	if _, err := rd.dc.DeleteSync(rd.ctx, &protodb.Entity{Key: key}); err != nil {
		return fmt.Errorf("remoteDB.DeleteSync: %w", err)
	}
	return nil
}

func (rd *RemoteDB) Set(key, value []byte) error {
	if _, err := rd.dc.Set(rd.ctx, &protodb.Entity{Key: key, Value: value}); err != nil {
		return fmt.Errorf("remoteDB.Set: %w", err)
	}
	return nil
}

func (rd *RemoteDB) SetSync(key, value []byte) error {
	if _, err := rd.dc.SetSync(rd.ctx, &protodb.Entity{Key: key, Value: value}); err != nil {
		return fmt.Errorf("remoteDB.SetSync: %w", err)
	}
	return nil
}

func (rd *RemoteDB) Get(key []byte) ([]byte, error) {
	res, err := rd.dc.Get(rd.ctx, &protodb.Entity{Key: key})
	if err != nil {
		return nil, fmt.Errorf("remoteDB.Get error: %w", err)
	}
	return res.Value, nil
}

func (rd *RemoteDB) Has(key []byte) (bool, error) {
	res, err := rd.dc.Has(rd.ctx, &protodb.Entity{Key: key})
	if err != nil {
		return false, err
	}
	return res.Exists, nil
}

func (rd *RemoteDB) NewBatch() tmdb.Batch {
	return newBatch(rd)
}

// TODO: Implement Print when tmdb.DB implements a method
// to print to a string and not db.Print to stdout.
func (rd *RemoteDB) Print() error {
	return errors.New("remoteDB.Print: unimplemented")
}

func (rd *RemoteDB) Stats() map[string]string {
	stats, err := rd.dc.Stats(rd.ctx, &protodb.Nothing{})
	if err != nil || stats == nil {
		return nil
	}
	return stats.Data
}

func (rd *RemoteDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	dic, err := rd.dc.Iterator(rd.ctx, &protodb.Entity{Start: start, End: end})
	if err != nil {
		return nil, fmt.Errorf("RemoteDB.Iterator error: %w", err)
	}
	return makeIterator(dic), nil
}

func (rd *RemoteDB) PrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	return rd.Iterator(start, end)
}

func (rd *RemoteDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	dic, err := rd.dc.ReverseIterator(rd.ctx, &protodb.Entity{Start: start, End: end})
	if err != nil {
		return nil, fmt.Errorf("RemoteDB.Iterator error: %w", err)
	}
	return makeReverseIterator(dic), nil
}

func (rd *RemoteDB) ReversePrefixIterator(prefix []byte) (tmdb.Iterator, error) {
	start, end, err := util.PrefixToRange(prefix)
	if err != nil {
		return nil, err
	}
	return rd.ReverseIterator(start, end)
}
