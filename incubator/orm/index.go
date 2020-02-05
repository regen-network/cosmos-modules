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

// indexer creates and modifies the second MultiKeyIndex based on the operations and changes on the primary object.
type indexer interface {
	OnCreate(store sdk.KVStore, rowID uint64, value interface{}) error
	OnDelete(store sdk.KVStore, rowID uint64, value interface{}) error
	OnUpdate(store sdk.KVStore, rowID uint64, newValue, oldValue interface{}) error
}

// MultiKeyIndex is an index where multiple entries can point to the same underlying object as opposite to a unique index
// where only one entry is allowed.
type MultiKeyIndex struct {
	storeKey  sdk.StoreKey
	prefix    byte
	rowGetter func(ctx HasKVStore, rowID uint64, dest interface{}) (key []byte, err error)
	indexer   indexer
}

func NewIndex(builder Indexable, prefix byte, indexer IndexerFunc) *MultiKeyIndex {
	idx := MultiKeyIndex{
		storeKey:  builder.StoreKey(),
		prefix:    prefix,
		rowGetter: builder.RowGetter(),
		indexer:   NewIndexer(indexer),
	}
	builder.AddAfterSaveInterceptor(idx.onSave)
	builder.AddAfterDeleteInterceptor(idx.onDelete)
	return &idx
}

func (i MultiKeyIndex) Has(ctx HasKVStore, key []byte) bool {
	// can only be answered by a prefix scan. see makeIndexPrefixScanKey
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	defer it.Close()
	return it.Valid()
}

func (i MultiKeyIndex) Get(ctx HasKVStore, key []byte) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil
}

// PrefixScan returns an Iterator over a domain of keys in ascending order. End is exclusive.
// Start is an MultiKeyIndex key or prefix. It must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
// To iterate over entire domain, use PrefixScan(nil, nil)
//
// WARNING: The use of a PrefixScan can be very expensive in terms of Gas. Please make sure you do not expose
// this as an endpoint to the public without further limits.
// Example:
//			it, err := idx.PrefixScan(ctx, start, end)
//			if err !=nil {
//				return err
//			}
//			const defaultLimit = 20
//			it = LimitIterator(it, defaultLimit)
//
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
func (i MultiKeyIndex) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	if start != nil && end != nil && bytes.Compare(start, end) >= 0 {
		return NewInvalidIterator(), errors.Wrap(ErrArgument, "start must be less than end")
	}
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	it := store.Iterator(start, end)
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil
}

// ReversePrefixScan returns an Iterator over a domain of keys in descending order. End is exclusive.
// Start is an MultiKeyIndex key or prefix. It must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
// To iterate over entire domain, use PrefixScan(nil, nil)
//
// WARNING: The use of a ReversePrefixScan can be very expensive in terms of Gas. Please make sure you do not expose
// this as an endpoint to the public without further limits. See `LimitIterator`
//
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
func (i MultiKeyIndex) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	if start != nil && end != nil && bytes.Compare(start, end) >= 0 {
		return NewInvalidIterator(), errors.Wrap(ErrArgument, "start must be less than end")
	}
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	it := store.ReverseIterator(start, end)
	return indexIterator{ctx: ctx, it: it, rowGetter: i.rowGetter}, nil
}

func (i MultiKeyIndex) onSave(ctx HasKVStore, rowID uint64, newValue, oldValue interface{}) error {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	if oldValue == nil {
		return i.indexer.OnCreate(store, rowID, newValue)
	}
	return i.indexer.OnUpdate(store, rowID, newValue, oldValue)
}

func (i MultiKeyIndex) onDelete(ctx HasKVStore, rowID uint64, oldValue interface{}) error {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	return i.indexer.OnDelete(store, rowID, oldValue)
}

type UniqueIndex struct {
	MultiKeyIndex
}

// NewUniqueIndex create a new Index object where duplicate keys are prohibited.
func NewUniqueIndex(builder Indexable, prefix byte, uniqueIndexerFunc UniqueIndexerFunc) *UniqueIndex {
	idx := UniqueIndex{
		MultiKeyIndex: MultiKeyIndex{
			storeKey:  builder.StoreKey(),
			prefix:    prefix,
			rowGetter: builder.RowGetter(),
			indexer:   NewUniqueIndexer(uniqueIndexerFunc),
		},
	}
	builder.AddAfterSaveInterceptor(idx.onSave)
	builder.AddAfterDeleteInterceptor(idx.onDelete)
	return &idx
}

// RowID looks up the rowID for an MultiKeyIndex key. Returns `ErrNotFound` when not exists in MultiKeyIndex.
func (i UniqueIndex) RowID(ctx HasKVStore, key []byte) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(i.storeKey), []byte{i.prefix})
	it := store.Iterator(makeIndexPrefixScanKey(key, 0), makeIndexPrefixScanKey(key, math.MaxUint64))
	defer it.Close()
	if !it.Valid() {
		return 0, ErrNotFound
	}
	return stripRowIDFromIndexPrefixScanKey(it.Key()), nil
}

// indexIterator uses rowGetter to lazy load new model values on request.
type indexIterator struct {
	ctx       HasKVStore
	rowGetter RowGetter
	it        types.Iterator
}

// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
// are no more items the ErrIteratorDone error is returned
// The key is the rowID and not any MultiKeyIndex key.
func (i indexIterator) LoadNext(dest interface{}) ([]byte, error) {
	if !i.it.Valid() {
		return nil, ErrIteratorDone
	}
	indexPrefixKey := i.it.Key()
	rowID := stripRowIDFromIndexPrefixScanKey(indexPrefixKey)
	i.it.Next()
	return i.rowGetter(i.ctx, rowID, dest)
}

// Close releases the iterator and should be called at the end of iteration
func (i indexIterator) Close() error {
	i.it.Close()
	return nil
}

func stripRowIDFromIndexPrefixScanKey(indexPrefixKey []byte) uint64 {
	n := len(indexPrefixKey)
	return binary.BigEndian.Uint64(indexPrefixKey[n-8:])
}
