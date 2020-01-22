package orm

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ Indexable = &AutoUInt64TableBuilder{}

func NewAutoUInt64TableBuilder(prefixData byte, prefixSeq byte, key sdk.StoreKey, cdc *codec.Codec, model interface{}) *AutoUInt64TableBuilder {
	if prefixData == prefixSeq {
		panic("prefixData and prefixSeq must be unique")
	}
	if cdc == nil {
		panic("codec must not be empty")
	}
	if model == nil {
		panic("model must not be empty")
	}
	tp := reflect.TypeOf(model)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	return &AutoUInt64TableBuilder{prefixData: prefixData, prefixSeq: prefixSeq, storeKey: key, cdc: cdc, model: tp}
}

type AutoUInt64TableBuilder struct {
	model       reflect.Type
	prefixData  byte
	prefixSeq   byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

// todo: this function gives access to the storage. It does not really fit the builder patter.
func (a AutoUInt64TableBuilder) RowGetter() RowGetter {
	return NewTypeSafeRowGetter(a.storeKey, a.prefixData, a.cdc, a.model)
}

func (a AutoUInt64TableBuilder) StoreKey() sdk.StoreKey {
	return a.storeKey
}

func (a AutoUInt64TableBuilder) Build() AutoUInt64Table {
	seq := NewSequence(a.storeKey, a.prefixSeq)
	return AutoUInt64Table{
		model:       a.model,
		sequence:    seq,
		prefix:      a.prefixData,
		storeKey:    a.storeKey,
		cdc:         a.cdc,
		afterSave:   a.afterSave,
		afterDelete: a.afterDelete,
	}
}

func (a *AutoUInt64TableBuilder) AddAfterSaveInterceptor(interceptor AfterSaveInterceptor) {
	a.afterSave = append(a.afterSave, interceptor)
}

func (a *AutoUInt64TableBuilder) AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor) {
	a.afterDelete = append(a.afterDelete, interceptor)
}

type AutoUInt64Table struct {
	model       reflect.Type
	prefix      byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	sequence    *Sequence
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

func (a AutoUInt64Table) Create(ctx HasKVStore, obj interface{}) (uint64, error) {
	if err := assertCorrectType(a.model, obj); err != nil {
		return 0, err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	v, err := a.cdc.MarshalBinaryBare(obj)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to serialize %T", obj)
	}
	rowID := a.sequence.NextVal(ctx)

	key := EncodeSequence(rowID)
	store.Set(key, v)
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, obj, nil); err != nil {
			return 0, errors.Wrapf(err, "interceptor %d failed", i)
		}
	}

	return rowID, nil
}

func (a AutoUInt64Table) Save(ctx HasKVStore, rowID uint64, newValue interface{}) error {
	if err := assertCorrectType(a.model, newValue); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	var oldValue = reflect.New(a.model).Interface()
	_, err := a.GetOne(ctx, rowID, oldValue)
	if err != nil {
		return errors.Wrap(err, "load old value")
	}
	newValueEncoded, err := a.cdc.MarshalBinaryBare(newValue)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize %T", newValue)
	}

	key := EncodeSequence(rowID)
	store.Set(key, newValueEncoded)
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, newValue, oldValue); err != nil {
			return errors.Wrapf(err, "interceptor %d failed", i)
		}
	}
	return nil
}

func (a AutoUInt64Table) Delete(ctx HasKVStore, rowID uint64) error {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	key := EncodeSequence(rowID)

	var oldValue = reflect.New(a.model).Interface()
	_, err := a.GetOne(ctx, rowID, oldValue)
	if err != nil {
		return errors.Wrap(err, "load old value")
	}
	store.Delete(key)

	for i, itc := range a.afterDelete {
		if err := itc(ctx, rowID, oldValue); err != nil {
			return errors.Wrapf(err, "delete interceptor %d failed", i)
		}
	}
	return nil
}

func (a AutoUInt64Table) Has(ctx HasKVStore, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	return store.Has(EncodeSequence(id))
}

func (a AutoUInt64Table) GetOne(ctx HasKVStore, rowID uint64, dest interface{}) ([]byte, error) {
	x := NewTypeSafeRowGetter(a.storeKey, a.prefix, a.cdc, a.model)
	return x(ctx, rowID, dest)
}

// end is not included
func (a AutoUInt64Table) PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	if start >= end {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	return &autoUInt64Iterator{
		ctx:       ctx,
		rowGetter: NewTypeSafeRowGetter(a.storeKey, a.prefix, a.cdc, a.model),
		it:        store.Iterator(EncodeSequence(start), EncodeSequence(end)),
	}, nil
}

func (a AutoUInt64Table) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	if start >= end {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	return &autoUInt64Iterator{
		ctx:       ctx,
		rowGetter: NewTypeSafeRowGetter(a.storeKey, a.prefix, a.cdc, a.model),
		it:        store.ReverseIterator(EncodeSequence(start), EncodeSequence(end)),
	}, nil
}

// autoUInt64Iterator
type autoUInt64Iterator struct {
	ctx       HasKVStore
	rowGetter RowGetter
	it        types.Iterator
}

func (i autoUInt64Iterator) LoadNext(dest interface{}) ([]byte, error) {
	if !i.it.Valid() {
		return nil, ErrIteratorDone
	}
	rowID := i.it.Key()
	i.it.Next()
	return i.rowGetter(i.ctx, DecodeSequence(rowID), dest)
}

func (i autoUInt64Iterator) Close() error {
	i.it.Close()
	return nil
}
