package orm

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewIndex(builder TableBuilder, prefix []byte, indexer IndexerFunc) Index {
	idx := index{
		prefix: prefix,
	}
	//builder.AddAfterDeleteInterceptor(idx.OnDelete)
	//builder.AddAfterSaveInterceptor(idx.OnSave)

	return &idx
}

// indexRef
type indexRef struct {
	prefix  []byte
	indexer Indexer
}

var _ Indexer = IndexerFunc(nil)

type IndexerFunc func(value interface{}) ([]byte, error)

func (i IndexerFunc) DoIndex(store sdk.KVStore, rowId uint64, key []byte, value interface{}) error {
	key, err := i(value)
	if err != nil {
		return err
	}
	indexKey := makeIndexPrefixScanKey(key, rowId)
	if !store.Has(indexKey) {
		store.Set(indexKey, []byte{0})
	}
	return nil
}

func (i IndexerFunc) BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index {
	return index{storeKey: storeKey, prefix: prefix, modelGetter: modelGetter}
}

func makeIndexPrefixScanKey(indexKey []byte, rowId uint64) []byte {
	n := len(indexKey)
	res := make([]byte, n+8)
	copy(res, indexKey)
	binary.BigEndian.PutUint64(res[n:], rowId)
	return res
}

type index struct {
	storeKey    sdk.StoreKey
	prefix      []byte
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
}

func (i index) Has(ctx HasKVStore, key []byte) (bool, error) {
	panic("implement me")
}

func (i index) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(key, nil)
	return indexIterator{ctx: ctx, it: it, end: key, modelGetter: i.modelGetter}, nil
}

func (i index) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (i index) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

// TODO: method signature does not make sense returning an error but the store panics
// id unused
func (i index) OnDelete(ctx HasKVStore, id uint64, key []byte) error {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	store.Delete(key)
	return nil
}

// TODO: method signature does not make sense returning an error but the store panics
// id unused
func (i index) OnSave(ctx HasKVStore, id uint64, key []byte, value interface{}) error {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	store.Set(key, []byte{})
	return nil
}

type indexIterator struct {
	ctx         HasKVStore
	modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
	it          types.Iterator
	end         []byte
	reverse     bool
}

func (i indexIterator) LoadNext(dest interface{}) (key []byte, err error) {
	if !i.it.Valid() {
		return nil, fmt.Errorf("not found")
	}
	indexPrefixKey := i.it.Key()
	n := len(indexPrefixKey)
	indexKey := indexPrefixKey[:n-8]
	cmp := bytes.Compare(indexKey, i.end)
	if i.end != nil {
		if !i.reverse && cmp > 0 {
			return nil, fmt.Errorf("not found")
		} else if i.reverse && cmp < 0 {
			return nil, fmt.Errorf("not found")
		}
	}
	rowId := binary.BigEndian.Uint64(indexPrefixKey[n-8:])
	i.it.Next()
	return i.modelGetter(i.ctx, rowId, dest)
}

func (i indexIterator) Close() error {
	i.it.Close()
	return nil
}
