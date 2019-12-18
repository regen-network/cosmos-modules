package orm

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


// IndexerFunc creates or or multiple index keyrs from the source object.
type IndexerFunc func(value interface{}) ([][]byte, error)

// todo: store is a prefix store. type should be explicit so people know. maybe make this a private method instead?
func (i IndexerFunc) OnCreate(store sdk.KVStore, rowId uint64, value interface{}) error {
	secondaryIndexKeys, err := i(value) // todo: multiple index keys?
	if err != nil {
		return err
	}
	for _, secondaryIndexKey := range secondaryIndexKeys {
		indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
		if !store.Has(indexKey) {
			store.Set(indexKey, []byte{0})
		}
	}
	return nil
}

// todo: same comments as OnCreate
func (i IndexerFunc) OnDelete(store sdk.KVStore, rowId uint64, value interface{}) error {
	secondaryIndexKeys, err := i(value) // todo: multiple index keys?
	if err != nil {
		return err
	}
	for _, secondaryIndexKey := range secondaryIndexKeys {
		indexKey := makeIndexPrefixScanKey(secondaryIndexKey, rowId)
		store.Delete(indexKey)
	}
	return nil
}

// todo: same comments as OnCreate
func (i IndexerFunc) OnUpdate(store sdk.KVStore, rowId uint64, newValue, oldValue interface{}) error {
	oldSecIdxKeys, err := i(oldValue)
	if err != nil {
		return err
	}
	newSecIdxKeys, err := i(newValue)
	if err != nil {
		return err
	}
	for _, oldIdxKey := range difference(oldSecIdxKeys, newSecIdxKeys) {
		store.Delete(makeIndexPrefixScanKey(oldIdxKey, rowId))
	}
	for _, newIdxKey := range difference(newSecIdxKeys, oldSecIdxKeys) {
		prefixedKey := makeIndexPrefixScanKey(newIdxKey, rowId)
		if !store.Has(prefixedKey) {
			store.Set(prefixedKey, []byte{0})
		}
	}
	return nil
}

// difference returns the list of elements that are in a but not in b.
func difference(a [][]byte, b [][]byte) [][]byte {
	set := make(map[string]struct{}, len(b))
	for _, v := range b {
		set[string(v)] = struct{}{}
	}
	var result [][]byte
	for _, v := range a {
		if _, ok := set[string(v)]; !ok {
			result = append(result, v)
		}
	}
	return result
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
	storeKey  sdk.StoreKey
	prefix    []byte
	rowGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
	indexer   IndexerFunc
}

func NewIndex(builder TableBuilder, prefix []byte, indexer IndexerFunc) *index {
	idx := index{
		storeKey:  builder.StoreKey(),
		prefix:    prefix,
		rowGetter: builder.RowGetter(),
		indexer:   indexer,
	}
	builder.AddAfterSaveInterceptor(idx.onSave)
	builder.AddAfterDeleteInterceptor(idx.onDelete)
	return &idx
}

// todo: store panics on errors. why return an error here?
func (i index) Has(ctx HasKVStore, key []byte) (bool, error) {
	//todo: does not work: return store.Has(key), nil
	// can only be answered by a prefix scan. see makeIndexPrefixScanKey

	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	defer it.Close()
	return it.Valid(), nil
}

func (i index) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	return indexIterator{ctx: ctx, it: it, end: key, modelGetter: i.rowGetter}, nil
}

// RowID looks up the rowID for an index key. Returns ErrNotFound when not exists in index.
func (i index) RowID(ctx HasKVStore, key []byte) (uint64, error) {
	// todo: support multiIndex ?
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	if !it.Valid() {
		return 0, ErrNotFound
	}
	return stripRowIDFromIndexKey(it.Key()), nil
}

func (i index) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (i index) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	panic("implement me")
}

func (i index) onSave(ctx HasKVStore, rowID uint64, key []byte, newValue, oldValue interface{}) error {
	// todo: this is the on create indexer, for update the old value may has to be removed
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	if oldValue == nil {
		return i.indexer.OnCreate(store, rowID, newValue)
	}
	return i.indexer.OnUpdate(store, rowID, newValue, oldValue)

}

func (i index) onDelete(ctx HasKVStore, rowId uint64, key []byte, oldValue interface{}) error {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	return i.indexer.OnDelete(store, rowId, oldValue)
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
		return nil, err
	}
	indexPrefixKey := i.it.Key()
	n := len(indexPrefixKey)
	indexKey := indexPrefixKey[:n-8]
	cmp := bytes.Compare(indexKey, i.end)
	if i.end != nil {
		if !i.reverse && cmp > 0 {
			return nil, err
		} else if i.reverse && cmp < 0 {
			return nil, err
		}
	}
	rowId := stripRowIDFromIndexKey(indexPrefixKey)
	i.it.Next()
	return i.modelGetter(i.ctx, rowId, dest)
}

func (i indexIterator) Close() error {
	i.it.Close()
	return nil
}

func stripRowIDFromIndexKey(indexPrefixKey []byte) uint64 {
	n := len(indexPrefixKey)
	return binary.BigEndian.Uint64(indexPrefixKey[n-8:])
}
