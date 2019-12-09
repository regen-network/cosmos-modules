package orm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ Indexer = IndexerFunc(nil)

type IndexerFunc func(value interface{}) ([]byte, error)

// todo: primary key is unused. same as rowID anyway
// todo: store is a prefix store. type should be explicit so people know. maybe make this a private method instead?
func (i IndexerFunc) DoIndex(store sdk.KVStore, rowId uint64, primaryKey []byte, value interface{}) error {
	secondaryIndexKey, err := i(value)
	if err != nil {
		return err
	}
	indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
	if !store.Has(indexKey) {
		store.Set(indexKey, []byte{0})
	}
	return nil
}

// TODO: remove function. there should only be 1 way to create an indexer: NewIndex
func (i IndexerFunc) BuildIndex(storeKey sdk.StoreKey, prefix []byte, modelGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)) Index {
	panic("what should we do here?")
}

// todo: this feels quite complicated when reading the data. Why not store rowID(s) as payload instead?
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
	indexer     IndexerFunc
}

func NewIndex(builder TableBuilder, prefix []byte, indexer IndexerFunc) Index {
	idx := index{
		storeKey:    builder.StoreKey(),
		prefix:      prefix,
		modelGetter: builder.ModelGetter(),
		indexer:     indexer,
	}
	builder.AddAfterSaveInterceptor(idx.onSave)
	return &idx
}

// todo: store panics on errors. why return an error here?
func (i index) Has(ctx HasKVStore, key []byte) (bool, error) {
	//todo: does not work: return store.Has(key), nil
	// can only be answered by a prefix scan. see makeIndexPrefixScanKey

	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	defer it.Close()
	println("+++ ", it.Valid())
	return it.Valid(), nil
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

func (i index) onSave(ctx HasKVStore, rowID uint64, key []byte, value interface{}) error {
	// todo: this is the on create indexer, for update the old value may has to be removed
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	err := i.indexer.DoIndex(store, rowID, key, value)
	if err != nil {
		return errors.Wrapf(err, "indexer for prefix %X failed", i.prefix)
	}
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
