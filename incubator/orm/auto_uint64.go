package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotFound     = errors.Register(errors.RootCodespace, 100, "not found")
	ErrIteratorDone = errors.Register(errors.RootCodespace, 101, "iterator done")
)

var _ TableBuilder = &autoUInt64TableBuilder{}

func NewAutoUInt64TableBuilder(prefix []byte, key sdk.StoreKey, cdc *codec.Codec) *autoUInt64TableBuilder {
	if len(prefix) == 0 {
		panic("prefix must not be empty")
	}
	if cdc == nil {
		panic("codec must not be empty")
	}
	return &autoUInt64TableBuilder{prefix: prefix, storeKey: key, cdc: cdc}
}

type autoUInt64TableBuilder struct {
	prefix      []byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	indexerRefs []indexRef
}

// todo: this function gives access to the storage. It does not really fit the builder patter.
func (a autoUInt64TableBuilder) ModelGetter() ModelGetter {
	return func(ctx HasKVStore, rowId uint64, dest interface{}) ([]byte, error) {
		store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
		key := encodeSequence(rowId)
		val := store.Get(key)
		// todo: how to handle not found?
		if val == nil {
			return nil, ErrNotFound // todo: discuss how to handle this scenario if we drop error return parameter
		}
		return key, a.cdc.UnmarshalBinaryBare(val, dest)
	}
}

func (a autoUInt64TableBuilder) StoreKey() sdk.StoreKey {
	return a.storeKey
}

func (a autoUInt64TableBuilder) RowGetter() RowGetter {
	panic("implement me")
}

// RegisterIndexer panics when prefix is a duplicate
func (a *autoUInt64TableBuilder) RegisterIndexer(prefix []byte, indexer Indexer) {
	// todo: fail on duplicates
	a.indexerRefs = append(a.indexerRefs, indexRef{prefix: prefix, indexer: indexer})
}

func (a autoUInt64TableBuilder) Build() autoUInt64Table {
	seq := NewSequence(a.storeKey, a.prefix)
	return autoUInt64Table{
		sequence:    seq,
		prefix:      a.prefix,
		storeKey:    a.storeKey,
		cdc:         a.cdc,
		indexerRefs: a.indexerRefs, // todo: clone slice
	}
}

// todo: required?
//func (a autoUInt64TableBuilder) AddAfterSaveInterceptor(interceptor AfterSaveInterceptor) {
//	panic("implement me")
//}
//
//func (a autoUInt64TableBuilder) AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor) {
//	panic("implement me")
//}

var _ AutoUInt64Table = autoUInt64Table{}

type autoUInt64Table struct {
	prefix      []byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	indexerRefs []indexRef
	sequence    Sequence
}

func (a autoUInt64Table) Create(ctx HasKVStore, value interface{}) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	v, err := a.cdc.MarshalBinaryBare(value)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to serialize %T", value)
	}
	id, err := a.sequence.NextVal(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "can not fetch next sequence value")
	}

	// todo: store does not return an error that we can handle or return
	key := encodeSequence(id)
	store.Set(key, v)
	for _, idx := range a.indexerRefs {
		store := prefix.NewStore(ctx.KVStore(a.storeKey), idx.prefix)
		err := idx.indexer.DoIndex(store, id, key, value)
		if err != nil {
			return 0, errors.Wrapf(err, "indexer for prefix %X failed", idx.prefix)
		}
	}
	return id, nil
}

func (a autoUInt64Table) Save(ctx HasKVStore, id uint64, value interface{}) error {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	v, err := a.cdc.MarshalBinaryBare(value)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize %T", value)
	}
	// todo: store does not return an error that we can handle or return
	store.Set(encodeSequence(id), v)
	return nil
}

// todo: there is no error result as store would panic
func (a autoUInt64Table) Has(ctx HasKVStore, id uint64) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	return store.Has(encodeSequence(id)), nil
}

func (a autoUInt64Table) Get(ctx HasKVStore, id uint64) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	key := encodeSequence(id)
	val := store.Get(key)
	if val == nil {
		return nil, ErrNotFound // todo: discuss how to handle this scenario if we drop error return parameter
	}
	return NewSingleValueIterator(a.cdc, key, val), nil // todo: SingleValueIterator is only used to satisfy the interface
}

func (a autoUInt64Table) PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

func (a autoUInt64Table) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	panic("implement me")
}

type iteratorFunc func(dest interface{}) (key []byte, err error)

func (i iteratorFunc) LoadNext(dest interface{}) (key []byte, err error) {
	return i(dest)
}

func (i iteratorFunc) Close() error {
	return nil
}

func NewSingleValueIterator(cdc *codec.Codec, key []byte, val []byte) Iterator {
	var closed bool
	return iteratorFunc(func(dest interface{}) ([]byte, error) {
		if closed || val == nil {
			return nil, ErrIteratorDone
		}
		closed = true
		return key, cdc.UnmarshalBinaryBare(val, dest)
	})
}