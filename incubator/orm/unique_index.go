package orm

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type uniqueIndexer struct {
	indexFn func(key []byte, value interface{}) ([]byte, error)
}

type uniqueIndex struct {
	uniqueIndexer
	storeKey    sdk.StoreKey
	prefix      []byte
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
}

func (u uniqueIndexer) DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error {
	indexKey, err := u.indexFn(key, value)
	if err != nil {
		return err
	}
	bz := store.Get(indexKey)
	if bz != nil {
		existingRowId := binary.BigEndian.Uint64(bz)
		if existingRowId != rowId {
			return fmt.Errorf("unique index constraint violated")
		}
		return nil
	}
	bz = make([]byte, 8)
	binary.BigEndian.PutUint64(bz, rowId)
	store.Set(indexKey, bz)
	return nil
}

func (u uniqueIndexer) BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index {
	return uniqueIndex{u, storeKey, prefix, modelGetter}
}

func (u uniqueIndex) Has(ctx HasKVStore, key []byte) (bool, error) {
	panic("implement me")
}

func (u uniqueIndex) getRowId(ctx HasKVStore, key []byte) uint64 {
	store := prefix.NewStore(ctx.KVStore(u.storeKey), u.prefix)
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (u uniqueIndex) GetOne(ctx HasKVStore, indexKey []byte, dest interface{}) ([]byte, error) {
	rowId := u.getRowId(ctx, indexKey)
	return u.modelGetter(ctx, rowId, dest)
}

func (u uniqueIndex) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	panic("implement me")
}

func (u uniqueIndex) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (u uniqueIndex) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func newPrimaryKeyIndexer() Indexer {
	return uniqueIndexer{indexFn: func(key []byte, value interface{}) (i []byte, e error) {
		return key, nil
	}}
}

func NewUniqueIndexer(fn IndexerFunc) Indexer {
	return uniqueIndexer{indexFn: func(key []byte, value interface{}) (i []byte, e error) {
		return fn(value)
	}}
}

var _ Indexer = uniqueIndexer{}

