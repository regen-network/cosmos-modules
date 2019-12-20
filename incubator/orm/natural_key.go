package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ TableBuilder = &naturalKeyTableBuilder{}

type naturalKeyer func(val interface{}) []byte // todo: note: in the api design this does not return an error unlike other indexer functions do

func NewNaturalKeyTableBuilder(prefixData, prefixSeq, prefixIndex []byte, key sdk.StoreKey, cdc *codec.Codec, model interface{}, getPrimaryKey naturalKeyer) *naturalKeyTableBuilder {
	if len(prefixIndex) == 0 {
		panic("prefixIndex must not be empty")
	}

	builder := NewAutoUInt64TableBuilder(prefixData, prefixSeq, key, cdc, model)

	idx := NewUniqueIndex(builder, prefixIndex, func(value interface{}) (bytes [][]byte, err error) {
		return [][]byte{getPrimaryKey(value)}, nil
	})
	return &naturalKeyTableBuilder{
		naturalKeyIndex:        idx,
		autoUInt64TableBuilder: builder,
	}
}

type naturalKeyTableBuilder struct {
	*autoUInt64TableBuilder
	naturalKeyIndex RowIDAwareIndex
}

func (a naturalKeyTableBuilder) Build() naturalKeyTable {
	return naturalKeyTable{
		autoTable:       a.autoUInt64TableBuilder.Build(),
		naturalKeyIndex: a.naturalKeyIndex,
	}
}

var _ NaturalKeyTable = naturalKeyTable{}

type naturalKeyTable struct {
	getPrimaryKey   naturalKeyer
	autoTable       autoUInt64Table
	naturalKeyIndex RowIDAwareIndex
}

func (a naturalKeyTable) GetOne(ctx HasKVStore, primKey []byte, dest interface{}) ([]byte, error) {
	it, err := a.Get(ctx, primKey)
	if err != nil {
		return nil, err
	}
	return First(it, dest)
}

func (a naturalKeyTable) Create(ctx HasKVStore, obj HasID) error {
	_, err := a.autoTable.Create(ctx, obj)
	return err
}

func (a naturalKeyTable) Save(ctx HasKVStore, newValue HasID) error {
	rowID, err := a.naturalKeyIndex.RowID(ctx, newValue.ID())
	if err != nil {
		return err
	}
	return a.autoTable.Save(ctx, rowID, newValue)
}

func (a naturalKeyTable) Delete(ctx HasKVStore, obj HasID) error {
	rowID, err := a.naturalKeyIndex.RowID(ctx, obj.ID())
	if err != nil {
		return err
	}
	return a.autoTable.Delete(ctx, rowID)
}

// todo: there is no error result as store would panic
func (a naturalKeyTable) Has(ctx HasKVStore, primKey []byte) (bool, error) {
	rowID, err := a.naturalKeyIndex.RowID(ctx, primKey)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return a.autoTable.Has(ctx, rowID)
}

func (a naturalKeyTable) Get(ctx HasKVStore, primKey []byte) (Iterator, error) {
	rowID, err := a.naturalKeyIndex.RowID(ctx, primKey)
	if err != nil {
		return nil, err
	}
	return a.autoTable.Get(ctx, rowID)
}

func (a naturalKeyTable) PrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	panic("implement me")
}

func (a naturalKeyTable) ReversePrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	panic("implement me")
}
