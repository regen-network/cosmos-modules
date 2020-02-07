package orm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Indexable = &NaturalKeyTableBuilder{}

func NewNaturalKeyTableBuilder(prefixData byte, storeKey sdk.StoreKey, model NaturalKeyed) *NaturalKeyTableBuilder {
	return &NaturalKeyTableBuilder{
		TableBuilder: NewTableBuilder(prefixData, storeKey, model),
	}
}

type NaturalKeyTableBuilder struct {
	*TableBuilder
}

func (a NaturalKeyTableBuilder) Build() NaturalKeyTable {
	return NaturalKeyTable{table: a.TableBuilder.Build()}

}

// NaturalKeyed defines an object type that is aware of it's immutable natural key.
type NaturalKeyed interface {
	// NaturalKey returns the immutable and serialized natural key of this object
	NaturalKey() RowID
	Persistent
}

// NaturalKeyTable provides simpler object style orm methods without passing database rowIDs.
// Entries are persisted and loaded with a reference to their natural key.
type NaturalKeyTable struct {
	table Table
}

func (a NaturalKeyTable) Create(ctx HasKVStore, obj NaturalKeyed) error {
	rowID := obj.NaturalKey()
	if a.table.Has(ctx, rowID) {
		return ErrUniqueConstraint
	}
	return a.table.Create(ctx, rowID, obj)
}

func (a NaturalKeyTable) Save(ctx HasKVStore, newValue NaturalKeyed) error {
	return a.table.Save(ctx, newValue.NaturalKey(), newValue)
}

func (a NaturalKeyTable) Delete(ctx HasKVStore, obj NaturalKeyed) error {
	return a.table.Delete(ctx, obj.NaturalKey())
}

func (a NaturalKeyTable) Has(ctx HasKVStore, naturalKey RowID) bool {
	return a.table.Has(ctx, naturalKey)
}

func (a NaturalKeyTable) GetOne(ctx HasKVStore, primKey RowID, dest Persistent) error {
	return a.table.GetOne(ctx, primKey, dest)
}

func (a NaturalKeyTable) PrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	return a.table.PrefixScan(ctx, start, end)
}

func (a NaturalKeyTable) ReversePrefixScan(ctx HasKVStore, start, end []byte) (Iterator, error) {
	return a.table.ReversePrefixScan(ctx, start, end)
}
