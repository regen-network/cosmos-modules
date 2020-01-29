package orm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ Indexable = &NaturalKeyTableBuilder{}

func NewNaturalKeyTableBuilder(prefixData, prefixSeq, prefixIndex byte, key sdk.StoreKey, cdc *codec.Codec, model NaturalKeyed) *NaturalKeyTableBuilder {
	if prefixIndex == prefixData || prefixData == prefixSeq {
		panic("prefixIndex must be unique")
	}

	builder := NewAutoUInt64TableBuilder(prefixData, prefixSeq, key, cdc, model)

	idx := NewUniqueIndex(builder, prefixIndex, func(value interface{}) ([]byte, error) {
		obj, ok := value.(NaturalKeyed)
		if !ok {
			return nil, errors.Wrapf(ErrType, "%T", value)
		}
		return obj.NaturalKey(), nil
	})
	return &NaturalKeyTableBuilder{
		naturalKeyIndex:        idx,
		AutoUInt64TableBuilder: builder,
	}
}

type NaturalKeyTableBuilder struct {
	*AutoUInt64TableBuilder
	naturalKeyIndex *UniqueIndex
}

func (a NaturalKeyTableBuilder) Build() NaturalKeyTable {
	return NaturalKeyTable{
		autoTable:       a.AutoUInt64TableBuilder.Build(),
		naturalKeyIndex: a.naturalKeyIndex,
	}
}

// NaturalKeyed defines an object type that is aware of it's immutable natural key.
type NaturalKeyed interface {
	// NaturalKey returns the immutable and serialized natural key of this object
	NaturalKey() []byte
}

// NaturalKeyTable provides simpler object style orm methods without passing database rowIDs.
// Entries are persisted and loaded with a reference to their natural key.
type NaturalKeyTable struct {
	autoTable       AutoUInt64Table
	naturalKeyIndex *UniqueIndex
}

func (a NaturalKeyTable) Create(ctx HasKVStore, obj NaturalKeyed) error {
	_, err := a.autoTable.Create(ctx, obj)
	return err
}

func (a NaturalKeyTable) Save(ctx HasKVStore, newValue NaturalKeyed) error {
	rowID, err := a.naturalKeyIndex.RowID(ctx, newValue.NaturalKey())
	if err != nil {
		return err
	}
	return a.autoTable.Save(ctx, rowID, newValue)
}

func (a NaturalKeyTable) Delete(ctx HasKVStore, obj NaturalKeyed) error {
	rowID, err := a.naturalKeyIndex.RowID(ctx, obj.NaturalKey())
	if err != nil {
		return err
	}
	return a.autoTable.Delete(ctx, rowID)
}

func (a NaturalKeyTable) Has(ctx HasKVStore, primKey []byte) bool {
	rowID, err := a.naturalKeyIndex.RowID(ctx, primKey)
	if err != nil {
		if err == ErrNotFound {
			return false
		}
		return false
	}
	return a.autoTable.Has(ctx, rowID)
}

func (a NaturalKeyTable) GetOne(ctx HasKVStore, primKey []byte, dest interface{}) ([]byte, error) {
	rowID, err := a.naturalKeyIndex.RowID(ctx, primKey)
	if err != nil {
		return nil, err
	}
	return a.autoTable.GetOne(ctx, rowID, dest)
}

func (a NaturalKeyTable) PrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	return a.naturalKeyIndex.PrefixScan(ctx, start, end)
}

func (a NaturalKeyTable) ReversePrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	return a.naturalKeyIndex.ReversePrefixScan(ctx, start, end)
}
