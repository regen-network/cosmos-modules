package orm

import (
	"bytes"
	"reflect"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ Indexable = &TableBuilder{}

func NewTableBuilder(prefixData byte, key sdk.StoreKey, model Persistent) *TableBuilder {
	if model == nil {
		panic("model must not be nil")
	}
	tp := reflect.TypeOf(model)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	return &TableBuilder{
		prefixData: prefixData,
		storeKey:   key,
		model:      tp,
	}
}

type TableBuilder struct {
	model       reflect.Type
	prefixData  byte
	storeKey    sdk.StoreKey
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

func (a TableBuilder) RowGetter() RowGetter {
	return NewTypeSafeRowGetter(a.storeKey, a.prefixData, a.model)
}

func (a TableBuilder) StoreKey() sdk.StoreKey {
	return a.storeKey
}

func (a TableBuilder) Build() Table {
	return Table{
		model:       a.model,
		prefix:      a.prefixData,
		storeKey:    a.storeKey,
		afterSave:   a.afterSave,
		afterDelete: a.afterDelete,
	}
}
func (a *TableBuilder) AddAfterSaveInterceptor(interceptor AfterSaveInterceptor) {
	a.afterSave = append(a.afterSave, interceptor)
}

func (a *TableBuilder) AddAfterDeleteInterceptor(interceptor AfterDeleteInterceptor) {
	a.afterDelete = append(a.afterDelete, interceptor)
}

type Table struct {
	model       reflect.Type
	prefix      byte
	storeKey    sdk.StoreKey
	afterSave   []AfterSaveInterceptor
	afterDelete []AfterDeleteInterceptor
}

func (a Table) Create(ctx HasKVStore, rowID RowID, obj Persistent) error {
	if err := assertCorrectType(a.model, obj); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	v, err := obj.Marshal()
	if err != nil {
		return errors.Wrapf(err, "failed to serialize %T", obj)
	}
	store.Set(rowID, v)
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, obj, nil); err != nil {
			return errors.Wrapf(err, "interceptor %d failed", i)
		}
	}
	return nil
}

func (a Table) Save(ctx HasKVStore, rowID RowID, newValue Persistent) error {
	if err := assertCorrectType(a.model, newValue); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	var oldValue = reflect.New(a.model).Interface().(Persistent)

	if err := a.GetOne(ctx, rowID, oldValue); err != nil {
		return errors.Wrap(err, "load old value")
	}
	newValueEncoded, err := newValue.Marshal()
	if err != nil {
		return errors.Wrapf(err, "failed to serialize %T", newValue)
	}

	store.Set(rowID, newValueEncoded)
	for i, itc := range a.afterSave {
		if err := itc(ctx, rowID, newValue, oldValue); err != nil {
			return errors.Wrapf(err, "interceptor %d failed", i)
		}
	}
	return nil
}

func (a Table) Delete(ctx HasKVStore, rowID RowID) error {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})

	var oldValue = reflect.New(a.model).Interface().(Persistent)
	if err := a.GetOne(ctx, rowID, oldValue); err != nil {
		return errors.Wrap(err, "load old value")
	}
	store.Delete(rowID)

	for i, itc := range a.afterDelete {
		if err := itc(ctx, rowID, oldValue); err != nil {
			return errors.Wrapf(err, "delete interceptor %d failed", i)
		}
	}
	return nil
}

func (a Table) Has(ctx HasKVStore, rowID RowID) bool {
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	it := store.Iterator(prefixRange(rowID))
	defer it.Close()
	return it.Valid()
}

func (a Table) GetOne(ctx HasKVStore, rowID RowID, dest Persistent) error {
	x := NewTypeSafeRowGetter(a.storeKey, a.prefix, a.model)
	return x(ctx, rowID, dest)
}

func (a Table) PrefixScan(ctx HasKVStore, start, end RowID) (Iterator, error) {
	if start != nil && end != nil && bytes.Compare(start, end) >= 0 {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	return &typeSafeIterator{
		ctx:       ctx,
		rowGetter: NewTypeSafeRowGetter(a.storeKey, a.prefix, a.model),
		it:        store.Iterator(start, end),
	}, nil
}

func (a Table) ReversePrefixScan(ctx HasKVStore, start, end RowID) (Iterator, error) {
	if start != nil && end != nil && bytes.Compare(start, end) >= 0 {
		return nil, errors.Wrap(ErrArgument, "start must be before end")
	}
	store := prefix.NewStore(ctx.KVStore(a.storeKey), []byte{a.prefix})
	return &typeSafeIterator{
		ctx:       ctx,
		rowGetter: NewTypeSafeRowGetter(a.storeKey, a.prefix, a.model),
		it:        store.ReverseIterator(start, end),
	}, nil
}

type typeSafeIterator struct {
	ctx       HasKVStore
	rowGetter RowGetter
	it        types.Iterator
}

func (i typeSafeIterator) LoadNext(dest Persistent) (RowID, error) {
	if !i.it.Valid() {
		return nil, ErrIteratorDone
	}
	rowID := i.it.Key()
	i.it.Next()
	return rowID, i.rowGetter(i.ctx, rowID, dest)
}

func (i typeSafeIterator) Close() error {
	i.it.Close()
	return nil
}
