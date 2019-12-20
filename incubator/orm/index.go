package orm

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// indexer creates and modifies the second index based on the operations and changes on the primary object.
type indexer interface {
	OnCreate(store sdk.KVStore, rowId uint64, value interface{}) error
	OnDelete(store sdk.KVStore, rowId uint64, value interface{}) error
	OnUpdate(store sdk.KVStore, rowId uint64, newValue, oldValue interface{}) error
}

// index is the entry point for all index related operations.
type index struct {
	storeKey  sdk.StoreKey
	prefix    []byte
	rowGetter func(ctx HasKVStore, rowId uint64, dest interface{}) (key []byte, err error)
	indexer   indexer
}

func NewIndex(builder TableBuilder, prefix []byte, indexer IndexerFunc) *index {
	idx := index{
		storeKey:  builder.StoreKey(),
		prefix:    prefix,
		rowGetter: builder.RowGetter(),
		indexer:   NewIndexer(indexer, false),
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
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil
}

func (i index) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	if bytes.Compare(start, end) > 0 {
		return nil, errors.Wrap(ErrArgument, "start must not be greater than end")
	}
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(start, 0), makeIndexPrefixScanKey(end, math.MaxUint64))
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil
}

func (i index) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.ReverseIterator(makeIndexPrefixScanKey(start, 0), makeIndexPrefixScanKey(end, math.MaxUint64))
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil

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

type uniqueIndex struct {
	index
}

func NewUniqueIndex(builder TableBuilder, prefix []byte, indexerFunc IndexerFunc) *uniqueIndex {
	idx := uniqueIndex{
		index: index{
			storeKey:  builder.StoreKey(),
			prefix:    prefix,
			rowGetter: builder.RowGetter(),
			indexer:   NewIndexer(indexerFunc, true),
		},
	}
	builder.AddAfterSaveInterceptor(idx.onSave)
	builder.AddAfterDeleteInterceptor(idx.onDelete)
	return &idx
}

// RowID looks up the rowID for an index key. Returns ErrNotFound when not exists in index.
func (i uniqueIndex) RowID(ctx HasKVStore, key []byte) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), i.prefix)
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	defer it.Close()
	if !it.Valid() {
		return 0, ErrNotFound
	}
	return stripRowIDFromIndexKey(it.Key()), nil
}

// indexIterator uses rowGetter to lazy load new model values on request.
type indexIterator struct {
	ctx       HasKVStore
	rowGetter RowGetter
	it        types.Iterator
}

func (i indexIterator) LoadNext(dest interface{}) ([]byte, error) {
	if !i.it.Valid() {
		return nil, ErrIteratorDone
	}
	indexPrefixKey := i.it.Key()
	rowId := stripRowIDFromIndexKey(indexPrefixKey)
	i.it.Next()
	return i.rowGetter(i.ctx, rowId, dest)
}

func (i indexIterator) Close() error {
	i.it.Close()
	return nil
}

func stripRowIDFromIndexKey(indexPrefixKey []byte) uint64 {
	n := len(indexPrefixKey)
	return binary.BigEndian.Uint64(indexPrefixKey[n-8:])
}
