package orm

import (
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const ormCodespace = "orm"

var (
	// todo: ormCodespace ok or do we need to claim error codes somehow?
	ErrNotFound         = errors.Register(ormCodespace, 100, "not found")
	ErrIteratorDone     = errors.Register(ormCodespace, 101, "iterator done")
	ErrIteratorInvalid  = errors.Register(ormCodespace, 102, "iterator invalid")
	ErrType             = errors.Register(ormCodespace, 110, "invalid type")
	ErrUniqueConstraint = errors.Register(ormCodespace, 111, "unique constraint violation")
	ErrArgument         = errors.Register(ormCodespace, 112, "invalid argument")
)

var _ TableBuilder = &autoUInt64TableBuilder{}

func NewAutoUInt64TableBuilder(prefixData []byte, prefixSeq []byte, key sdk.StoreKey, cdc *codec.Codec, model interface{}) *autoUInt64TableBuilder {
	if len(prefixData) == 0 {
		panic("prefixData must not be empty")
	}
	if len(prefixSeq) == 0 {
		panic("prefixSeq must not be empty")
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
	return &autoUInt64TableBuilder{prefixData: prefixData, prefixSeq: prefixSeq, storeKey: key, cdc: cdc, model: tp}
}

type autoUInt64TableBuilder struct {
	model       reflect.Type
	prefixData  []byte
	prefixSeq   []byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

// todo: this function gives access to the storage. It does not really fit the builder patter.
func (a autoUInt64TableBuilder) RowGetter() RowGetter {
	return typeSafeRowGetter(a.storeKey, a.prefixData, a.cdc, a.model)
}

func typeSafeRowGetter(storeKey sdk.StoreKey, prefixKey []byte, cdc *codec.Codec, model reflect.Type) RowGetter {
	return func(ctx HasKVStore, rowId uint64, dest interface{}) ([]byte, error) {
		if err := assertCorrectType(model, dest); err != nil {
			return nil, err
		}
		store := prefix.NewStore(ctx.KVStore(storeKey), prefixKey)
		key := EncodeSequence(rowId)
		val := store.Get(key)
		// todo: how to handle not found?
		if val == nil {
			return nil, ErrNotFound // todo: discuss how to handle this scenario if we drop error return parameter
		}
		return key, cdc.UnmarshalBinaryBare(val, dest)
	}
}

func (a autoUInt64TableBuilder) StoreKey() sdk.StoreKey {
	return a.storeKey
}

func (a autoUInt64TableBuilder) Build() autoUInt64Table {
	seq := NewSequence(a.storeKey, a.prefixSeq)
	return autoUInt64Table{
		model:       a.model,
		sequence:    seq,
		prefix:      a.prefixData,
		storeKey:    a.storeKey,
		cdc:         a.cdc,
		afterSave:   a.afterSave,
		afterDelete: a.afterDelete,
	}
}

func (a *autoUInt64TableBuilder) AddAfterSaveInterceptor(interceptor AfterSaveInterceptor) {
	a.afterSave = append(a.afterSave, interceptor)
}

func (a *autoUInt64TableBuilder) AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor) {
	a.afterDelete = append(a.afterDelete, interceptor)
}

var _ AutoUInt64Table = autoUInt64Table{}

type autoUInt64Table struct {
	model       reflect.Type
	prefix      []byte
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	sequence    Sequence
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

func (a autoUInt64Table) Create(ctx HasKVStore, obj interface{}) (uint64, error) {
	if err := assertCorrectType(a.model, obj); err != nil {
		return 0, err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	v, err := a.cdc.MarshalBinaryBare(obj)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to serialize %T", obj)
	}
	rowID, err := a.sequence.NextVal(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "can not fetch next sequence value")
	}

	// todo: store does not return an error that we can handle or return
	key := EncodeSequence(rowID)
	store.Set(key, v)
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, key, obj, nil); err != nil {
			return 0, errors.Wrapf(err, "interceptor %d failed", i)
		}
	}

	return rowID, nil
}

func (a autoUInt64Table) Save(ctx HasKVStore, rowID uint64, newValue interface{}) error {
	if err := assertCorrectType(a.model, newValue); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	var oldValue = reflect.New(a.model).Interface()
	it, err := a.Get(ctx, rowID)
	if err != nil {
		return err
	}
	_, err = it.LoadNext(oldValue)
	if err != nil {
		return err
	}

	v, err := a.cdc.MarshalBinaryBare(newValue)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize %T", newValue)
	}
	// todo: store does not return an error that we can handle or return
	key := EncodeSequence(rowID)
	store.Set(key, v)
	// todo: impl interceptor calls
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, key, newValue, oldValue); err != nil {
			return errors.Wrapf(err, "interceptor %d failed", i)
		}
	}
	return nil
}

func (a autoUInt64Table) Delete(ctx HasKVStore, rowID uint64) error {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	key := EncodeSequence(rowID)

	var oldValue = reflect.New(a.model).Interface()
	it, err := a.Get(ctx, rowID)
	if err != nil {
		return err
	}
	_, err = it.LoadNext(oldValue)
	if err != nil {
		return err
	}
	store.Delete(key)

	for i, itc := range a.afterDelete {
		if err := itc(ctx, rowID, key, oldValue); err != nil {
			return errors.Wrapf(err, "delete interceptor %d failed", i)
		}
	}
	return nil
}

// todo: there is no error result as store would panic
func (a autoUInt64Table) Has(ctx HasKVStore, id uint64) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	return store.Has(EncodeSequence(id)), nil
}

func (a autoUInt64Table) Get(ctx HasKVStore, rowID uint64) (Iterator, error) {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	key := EncodeSequence(rowID)
	val := store.Get(key)
	if val == nil {
		return nil, ErrNotFound // todo: discuss how to handle this scenario if we drop error return parameter
	}
	return NewSingleValueIterator(a.cdc, key, val), nil // todo: SingleValueIterator is only used to satisfy the interface
}

// end is not included
func (a autoUInt64Table) PrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	if start >= end {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	return &autoUInt64Iterator{
		ctx:       ctx,
		rowGetter: typeSafeRowGetter(a.storeKey, a.prefix, a.cdc, a.model),
		it:        store.Iterator(EncodeSequence(start), EncodeSequence(end)),
	}, nil
}

func (a autoUInt64Table) ReversePrefixScan(ctx HasKVStore, start uint64, end uint64) (Iterator, error) {
	if start >= end {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), a.prefix)
	return &autoUInt64Iterator{
		ctx:       ctx,
		rowGetter: typeSafeRowGetter(a.storeKey, a.prefix, a.cdc, a.model),
		it:        store.ReverseIterator(EncodeSequence(start), EncodeSequence(end)),
	}, nil
}

// typesafe iterator
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

func assertCorrectType(model reflect.Type, obj interface{}) error {
	tp := reflect.TypeOf(obj)
	if tp.Kind() != reflect.Ptr {
		return errors.Wrap(ErrType, "model destination must be a pointer")
	}
	if model != tp.Elem() {
		return errors.Wrapf(ErrType, "can not use %T with this bucket", obj)
	}
	return nil
}
